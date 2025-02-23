/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/certwatcher"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/Tchoupinax/image-operator/graphql"
	"github.com/go-logr/logr"

	buildahiov1alpha1 "github.com/Tchoupinax/image-operator/api/buildah.io/v1alpha1"
	skopeoiov1alpha1 "github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	tchoupinaxiov1beta "github.com/Tchoupinax/image-operator/api/v1beta"

	"github.com/Tchoupinax/image-operator/internal/controller"
	buildahiocontroller "github.com/Tchoupinax/image-operator/internal/controller/buildah.io"
	corecontroller "github.com/Tchoupinax/image-operator/internal/controller/core"
	skopeocontroller "github.com/Tchoupinax/image-operator/internal/controller/skopeo.io"
	helpers "github.com/Tchoupinax/image-operator/internal/helpers"
	webhookcorev1 "github.com/Tchoupinax/image-operator/internal/webhook/v1"

	// +kubebuilder:scaffold:imports

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	prometheusReloadGauge = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "image_operator_reload_total",
			Help: "Number of reload proccessed",
		},
		[]string{"image"},
	)
	lastTimeImageWasReloaded = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "image_operator_last_time_image_was_reloaded",
			Help: "Timestamp of when the image was reloaded for the last time",
		},
		[]string{"image"},
	)
	lastTimeImagebuilderWasReloaded = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "image_operator_last_time_imagebuilder_was_reloaded",
			Help: "Timestamp of when the image builder was reloaded for the last time",
		},
		[]string{"imagebuilder"},
	)
	dockerhubQuotaLimit = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "image_operator_dockerhub_quota_limit",
			Help: "What is the quota limit (hard limit) for Docker.io API",
		},
		[]string{"ip"},
	)
	dockerhubQuotaRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "image_operator_dockerhub_quota_remaining",
			Help: "How many pull are available from Docker.io API",
		},
		[]string{"ip"},
	)
	imagebuilderBuildsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_operator_imagebuilder_build_count",
			Help: "Count of builds done by imagebuilder",
		},
	)
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(skopeoiov1alpha1.AddToScheme(scheme))
	utilruntime.Must(buildahiov1alpha1.AddToScheme(scheme))
	utilruntime.Must(tchoupinaxiov1beta.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme

	metrics.Registry.MustRegister(
		dockerhubQuotaLimit,
		dockerhubQuotaRemaining,
		imagebuilderBuildsCount,
		lastTimeImageWasReloaded,
		lastTimeImagebuilderWasReloaded,
		prometheusReloadGauge,
	)
}

func main() {
	if helpers.GetEnv("FEATURE_DOCKERHUB_RATE_LIMIT_ENABLED", "false") == "true" {
		go heartBeatDockerhub(setupLog)
	} else {
		// Set default value to prevent from returning 0
		dockerhubQuotaRemaining.WithLabelValues("0.0.0.0").Set(-1)
		dockerhubQuotaLimit.WithLabelValues("0.0.0.0").Set(-1)
	}

	var namespaces = strings.Split(helpers.GetEnv("FEATURE_COPY_ON_THE_FLY_NAMESPACES_ALLOWED", "*"), ",")

	cacheNamespaces := make(map[string]cache.Config)
	cacheOptions := cache.Options{
		DefaultNamespaces: cacheNamespaces,
	}
	if !(len(namespaces) == 1 && namespaces[0] == "*") {
		for _, namespace := range namespaces {
			cacheNamespaces[namespace] = cache.Config{}
		}
	}

	var metricsAddr string
	var metricsCertPath, metricsCertName, metricsCertKey string
	var webhookCertPath, webhookCertName, webhookCertKey string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	var tlsOpts []func(*tls.Config)
	flag.StringVar(&metricsAddr, "metrics-bind-address", "0", "The address the metrics endpoint binds to. "+
		"Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", true,
		"If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead.")
	flag.StringVar(&webhookCertPath, "webhook-cert-path", "", "The directory that contains the webhook certificate.")
	flag.StringVar(&webhookCertName, "webhook-cert-name", "tls.crt", "The name of the webhook certificate file.")
	flag.StringVar(&webhookCertKey, "webhook-cert-key", "tls.key", "The name of the webhook key file.")
	flag.StringVar(&metricsCertPath, "metrics-cert-path", "",
		"The directory that contains the metrics server certificate.")
	flag.StringVar(&metricsCertName, "metrics-cert-name", "tls.crt", "The name of the metrics server certificate file.")
	flag.StringVar(&metricsCertKey, "metrics-cert-key", "tls.key", "The name of the metrics server key file.")
	flag.BoolVar(&enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if os.Getenv("GRAPHQL_API_ENABLED") == "true" {
		graphql.StartGraphqlServer()
	}

	// Only start operator when it's in the cluster
	if os.Getenv("AVOID_OPERATOR") == "true" {
		fmt.Println("Operator won't start")
	} else {
		// if the enable-http2 flag is false (the default), http/2 should be disabled
		// due to its vulnerabilities. More specifically, disabling http/2 will
		// prevent from being vulnerable to the HTTP/2 Stream Cancellation and
		// Rapid Reset CVEs. For more information see:
		// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
		// - https://github.com/advisories/GHSA-4374-p667-p6c8
		disableHTTP2 := func(c *tls.Config) {
			setupLog.Info("disabling http/2")
			c.NextProtos = []string{"http/1.1"}
		}

		if !enableHTTP2 {
			tlsOpts = append(tlsOpts, disableHTTP2)
		}

		// Create watchers for metrics and webhooks certificates
		var metricsCertWatcher, webhookCertWatcher *certwatcher.CertWatcher

		// Initial webhook TLS options
		webhookTLSOpts := tlsOpts

		if len(webhookCertPath) > 0 {
			setupLog.Info("Initializing webhook certificate watcher using provided certificates",
				"webhook-cert-path", webhookCertPath, "webhook-cert-name", webhookCertName, "webhook-cert-key", webhookCertKey)

			var err error
			webhookCertWatcher, err = certwatcher.New(
				filepath.Join(webhookCertPath, webhookCertName),
				filepath.Join(webhookCertPath, webhookCertKey),
			)
			if err != nil {
				setupLog.Error(err, "Failed to initialize webhook certificate watcher")
				os.Exit(1)
			}

			webhookTLSOpts = append(webhookTLSOpts, func(config *tls.Config) {
				config.GetCertificate = webhookCertWatcher.GetCertificate
			})
		}

		webhookServer := webhook.NewServer(webhook.Options{
			TLSOpts: webhookTLSOpts,
		})

		// Metrics endpoint is enabled in 'config/default/kustomization.yaml'. The Metrics options configure the server.
		// More info:
		// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/metrics/server
		// - https://book.kubebuilder.io/reference/metrics.html
		metricsServerOptions := metricsserver.Options{
			BindAddress:   metricsAddr,
			SecureServing: secureMetrics,
			// TODO(user): TLSOpts is used to allow configuring the TLS config used for the server. If certificates are
			// not provided, self-signed certificates will be generated by default. This option is not recommended for
			// production environments as self-signed certificates do not offer the same level of trust and security
			// as certificates issued by a trusted Certificate Authority (CA). The primary risk is potentially allowing
			// unauthorized access to sensitive metrics data. Consider replacing with CertDir, CertName, and KeyName
			// to provide certificates, ensuring the server communicates using trusted and secure certificates.
			TLSOpts: tlsOpts,
		}

		if secureMetrics {
			// FilterProvider is used to protect the metrics endpoint with authn/authz.
			// These configurations ensure that only authorized users and service accounts
			// can access the metrics endpoint. The RBAC are configured in 'config/rbac/kustomization.yaml'. More info:
			// https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/metrics/filters#WithAuthenticationAndAuthorization
			metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization
		}

		// If the certificate is not specified, controller-runtime will automatically
		// generate self-signed certificates for the metrics server. While convenient for development and testing,
		// this setup is not recommended for production.
		//
		// TODO(user): If you enable certManager, uncomment the following lines:
		// - [METRICS-WITH-CERTS] at config/default/kustomization.yaml to generate and use certificates
		// managed by cert-manager for the metrics server.
		// - [PROMETHEUS-WITH-CERTS] at config/prometheus/kustomization.yaml for TLS certification.
		if len(metricsCertPath) > 0 {
			setupLog.Info("Initializing metrics certificate watcher using provided certificates",
				"metrics-cert-path", metricsCertPath, "metrics-cert-name", metricsCertName, "metrics-cert-key", metricsCertKey)

			var err error
			metricsCertWatcher, err = certwatcher.New(
				filepath.Join(metricsCertPath, metricsCertName),
				filepath.Join(metricsCertPath, metricsCertKey),
			)
			if err != nil {
				setupLog.Error(err, "to initialize metrics certificate watcher", "error", err)
				os.Exit(1)
			}

			metricsServerOptions.TLSOpts = append(metricsServerOptions.TLSOpts, func(config *tls.Config) {
				config.GetCertificate = metricsCertWatcher.GetCertificate
			})
		}

		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme: scheme,
			Client: client.Options{
				Cache: &client.CacheOptions{
					DisableFor: []client.Object{
						&corev1.Secret{},
						&corev1.ConfigMap{},
						&corev1.Pod{},
					},
				},
			},
			Metrics:                metricsServerOptions,
			WebhookServer:          webhookServer,
			HealthProbeBindAddress: probeAddr,
			LeaderElection:         enableLeaderElection,
			LeaderElectionID:       "dce26553.tchoupinax.dev",
			Cache:                  cacheOptions,
			// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
			// when the Manager ends. This requires the binary to immediately end when the
			// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
			// speeds up voluntary leader transitions as the new leader don't have to wait
			// LeaseDuration time first.
			//
			// In the default scaffold provided, the program ends immediately after
			// the manager stops, so would be fine to enable this option. However,
			// if you are doing or is intended to do any operation such as perform cleanups
			// after the manager stops then its usage might be unsafe.
			// LeaderElectionReleaseOnCancel: true,
		})
		if err != nil {
			setupLog.Error(err, "unable to start manager")
			os.Exit(1)
		}

		// LEGACY
		if err = (&skopeocontroller.LegacyImageReconciler{
			Client:                   mgr.GetClient(),
			Scheme:                   mgr.GetScheme(),
			PrometheusReloadGauge:    *prometheusReloadGauge,
			LastTimeImageWasReloaded: *lastTimeImageWasReloaded,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Image")
			os.Exit(1)
		}
		// LEGACY
		if err = (&buildahiocontroller.LegacyImageBuilderReconciler{
			Client:                  mgr.GetClient(),
			Scheme:                  mgr.GetScheme(),
			ImagebuilderBuildsCount: imagebuilderBuildsCount,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ImageBuilder")
			os.Exit(1)
		}

		if err = (&controller.ImageReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Image")
			os.Exit(1)
		}
		if err = (&controller.ImageBuilderReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ImageBuilder")
			os.Exit(1)
		}

		// nolint:goconst
		if os.Getenv("ENABLE_WEBHOOKS") != "false" {
			if err = webhookcorev1.SetupPodWebhookWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create webhook", "webhook", "Pod")
				os.Exit(1)
			}
		}

		// Activate copy on fly feature. Disabled by default
		if helpers.GetEnv("FEATURE_COPY_ON_THE_FLY_ENABLED", "false") == "true" {
			if err = (&corecontroller.PodReconciler{
				Client:                mgr.GetClient(),
				Scheme:                mgr.GetScheme(),
				OnFlyNamespaceAllowed: namespaces,
			}).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "Pod")
				os.Exit(1)
			}
		}

		// +kubebuilder:scaffold:builder

		if metricsCertWatcher != nil {
			setupLog.Info("Adding metrics certificate watcher to manager")
			if err := mgr.Add(metricsCertWatcher); err != nil {
				setupLog.Error(err, "unable to add metrics certificate watcher to manager")
				os.Exit(1)
			}
		}

		if webhookCertWatcher != nil {
			setupLog.Info("Adding webhook certificate watcher to manager")
			if err := mgr.Add(webhookCertWatcher); err != nil {
				setupLog.Error(err, "unable to add webhook certificate watcher to manager")
				os.Exit(1)
			}
		}

		if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
			setupLog.Error(err, "unable to set up health check")
			os.Exit(1)
		}
		if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
			setupLog.Error(err, "unable to set up ready check")
			os.Exit(1)
		}

		setupLog.Info("starting manager")
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}
	}
}

func heartBeatDockerhub(logger logr.Logger) {
	var frequencySecond = helpers.GetEnv("FEATURE_DOCKERHUB_RATE_LIMIT_FREQUENCY_SECOND", "60")
	value, err := strconv.Atoi(frequencySecond)
	if err != nil {
		logger.Error(err, "Value of FEATURE_DOCKERHUB_RATE_LIMIT_FREQUENCY_SECOND is not a correct integer")
		return
	}

	for range time.Tick(time.Second * time.Duration(value)) {
		result := helpers.GetDockerhubLimit(setupLog)
		if result.Succeeded {
			logger.Info(fmt.Sprintf("Dockerhub quota reminds %d/%d with %s", result.Remaining, result.Limit, result.Ip))
			dockerhubQuotaRemaining.WithLabelValues(result.Ip).Set(float64(result.Remaining))
			dockerhubQuotaLimit.WithLabelValues(result.Ip).Set(float64(result.Limit))
		} else {
			logger.Info("You are rate-limited by DockerHub")
		}
	}
}
