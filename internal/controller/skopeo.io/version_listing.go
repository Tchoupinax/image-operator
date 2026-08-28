package controller

import (
	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/go-logr/logr"
)

type listVersionsFunc func(
	logger logr.Logger,
	sourceName string,
	matchingString string,
	allowCandidateRelease bool,
	dockerhubAuth helpers.DockerHubAuth,
	awsPublicECR helpers.AWSPublicECR,
) []string

var listVersionsForImage listVersionsFunc = helpers.ListVersions
