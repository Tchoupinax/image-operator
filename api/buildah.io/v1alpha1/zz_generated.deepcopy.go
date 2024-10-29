//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageBuilder) DeepCopyInto(out *ImageBuilder) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageBuilder.
func (in *ImageBuilder) DeepCopy() *ImageBuilder {
	if in == nil {
		return nil
	}
	out := new(ImageBuilder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ImageBuilder) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageBuilderList) DeepCopyInto(out *ImageBuilderList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ImageBuilder, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageBuilderList.
func (in *ImageBuilderList) DeepCopy() *ImageBuilderList {
	if in == nil {
		return nil
	}
	out := new(ImageBuilderList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ImageBuilderList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageBuilderSpec) DeepCopyInto(out *ImageBuilderSpec) {
	*out = *in
	out.Image = in.Image
	out.Resources = in.Resources
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageBuilderSpec.
func (in *ImageBuilderSpec) DeepCopy() *ImageBuilderSpec {
	if in == nil {
		return nil
	}
	out := new(ImageBuilderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageBuilderStatus) DeepCopyInto(out *ImageBuilderStatus) {
	*out = *in
	if in.RanJobs != nil {
		in, out := &in.RanJobs, &out.RanJobs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageBuilderStatus.
func (in *ImageBuilderStatus) DeepCopy() *ImageBuilderStatus {
	if in == nil {
		return nil
	}
	out := new(ImageBuilderStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageEndpoint) DeepCopyInto(out *ImageEndpoint) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageEndpoint.
func (in *ImageEndpoint) DeepCopy() *ImageEndpoint {
	if in == nil {
		return nil
	}
	out := new(ImageEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resource) DeepCopyInto(out *Resource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resource.
func (in *Resource) DeepCopy() *Resource {
	if in == nil {
		return nil
	}
	out := new(Resource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourcesQuota) DeepCopyInto(out *ResourcesQuota) {
	*out = *in
	out.Limits = in.Limits
	out.Requests = in.Requests
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourcesQuota.
func (in *ResourcesQuota) DeepCopy() *ResourcesQuota {
	if in == nil {
		return nil
	}
	out := new(ResourcesQuota)
	in.DeepCopyInto(out)
	return out
}
