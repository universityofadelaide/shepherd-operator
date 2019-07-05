// +build !ignore_autogenerated

// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ingress.
func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Ingress) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressList) DeepCopyInto(out *IngressList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Ingress, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressList.
func (in *IngressList) DeepCopy() *IngressList {
	if in == nil {
		return nil
	}
	out := new(IngressList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IngressList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressSpec) DeepCopyInto(out *IngressSpec) {
	*out = *in
	in.Routes.DeepCopyInto(&out.Routes)
	in.Whitelist.DeepCopyInto(&out.Whitelist)
	out.Service = in.Service
	out.Prometheus = in.Prometheus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSpec.
func (in *IngressSpec) DeepCopy() *IngressSpec {
	if in == nil {
		return nil
	}
	out := new(IngressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressSpecPrometheus) DeepCopyInto(out *IngressSpecPrometheus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSpecPrometheus.
func (in *IngressSpecPrometheus) DeepCopy() *IngressSpecPrometheus {
	if in == nil {
		return nil
	}
	out := new(IngressSpecPrometheus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressSpecRoute) DeepCopyInto(out *IngressSpecRoute) {
	*out = *in
	if in.Subpaths != nil {
		in, out := &in.Subpaths, &out.Subpaths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSpecRoute.
func (in *IngressSpecRoute) DeepCopy() *IngressSpecRoute {
	if in == nil {
		return nil
	}
	out := new(IngressSpecRoute)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressSpecRoutes) DeepCopyInto(out *IngressSpecRoutes) {
	*out = *in
	in.Primary.DeepCopyInto(&out.Primary)
	if in.Secondary != nil {
		in, out := &in.Secondary, &out.Secondary
		*out = make([]IngressSpecRoute, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSpecRoutes.
func (in *IngressSpecRoutes) DeepCopy() *IngressSpecRoutes {
	if in == nil {
		return nil
	}
	out := new(IngressSpecRoutes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressSpecService) DeepCopyInto(out *IngressSpecService) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSpecService.
func (in *IngressSpecService) DeepCopy() *IngressSpecService {
	if in == nil {
		return nil
	}
	out := new(IngressSpecService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressStatus) DeepCopyInto(out *IngressStatus) {
	*out = *in
	out.CloudFront = in.CloudFront
	in.Certificate.DeepCopyInto(&out.Certificate)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressStatus.
func (in *IngressStatus) DeepCopy() *IngressStatus {
	if in == nil {
		return nil
	}
	out := new(IngressStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressStatusCertificateRef) DeepCopyInto(out *IngressStatusCertificateRef) {
	*out = *in
	in.Details.DeepCopyInto(&out.Details)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressStatusCertificateRef.
func (in *IngressStatusCertificateRef) DeepCopy() *IngressStatusCertificateRef {
	if in == nil {
		return nil
	}
	out := new(IngressStatusCertificateRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressStatusCloudFrontRef) DeepCopyInto(out *IngressStatusCloudFrontRef) {
	*out = *in
	out.Details = in.Details
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressStatusCloudFrontRef.
func (in *IngressStatusCloudFrontRef) DeepCopy() *IngressStatusCloudFrontRef {
	if in == nil {
		return nil
	}
	out := new(IngressStatusCloudFrontRef)
	in.DeepCopyInto(out)
	return out
}
