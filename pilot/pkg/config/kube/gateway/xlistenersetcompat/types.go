// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

// Package xlistenersetcompat provides backwards-compatibility types for the
// experimental gateway.networking.x-k8s.io/v1alpha1 XListenerSet kind. The
// type was removed from sigs.k8s.io/gateway-api in v1.5 when it graduated
// to GA as gateway.networking.k8s.io/v1.ListenerSet. Clusters that still
// have the experimental CRD installed (with attached HTTPRoutes parentRef-ing
// kind:XListenerSet) need pilot to ingest those objects during migration.
//
// The on-wire JSON shape is identical to the GA ListenerSet (verified by
// diffing gateway-api v1.4 XListenerSet types vs v1.5 ListenerSet types).
// Spec/Status types alias to the GA types so DeepCopy implementations
// delegate to upstream rather than maintaining a duplicate copy.
package xlistenersetcompat

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// GroupName / GroupVersion / SchemeGroupVersion identifiers for the
// experimental kind. These match the wire-level CRD registration.
const GroupName = "gateway.networking.x-k8s.io"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

// Type aliases to the GA gateway-api types. The on-wire JSON shape is
// identical so we delegate to upstream rather than duplicating definitions.
// Aliases let istio's schema codegen resolve `proto: "ListenerSetSpec"`
// against this package without needing a separate GA package import in
// the generated files.
type (
	ListenerSetSpec   = gatewayv1.ListenerSetSpec
	ListenerSetStatus = gatewayv1.ListenerSetStatus
)

// XListenerSet is the in-fork compat type for
// gateway.networking.x-k8s.io/v1alpha1.XListenerSet. Field types delegate
// to the GA gateway-api/apis/v1 package since the schemas are identical.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type XListenerSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   gatewayv1.ListenerSetSpec   `json:"spec"`
	Status gatewayv1.ListenerSetStatus `json:"status,omitempty"`
}

// XListenerSetList is the list type for XListenerSet.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type XListenerSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []XListenerSet `json:"items"`
}

// DeepCopyInto copies the receiver into out.
func (in *XListenerSet) DeepCopyInto(out *XListenerSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy returns a deep copy of the receiver.
func (in *XListenerSet) DeepCopy() *XListenerSet {
	if in == nil {
		return nil
	}
	out := new(XListenerSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject satisfies runtime.Object.
func (in *XListenerSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto copies the receiver into out.
func (in *XListenerSetList) DeepCopyInto(out *XListenerSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]XListenerSet, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

// DeepCopy returns a deep copy of the receiver.
func (in *XListenerSetList) DeepCopy() *XListenerSetList {
	if in == nil {
		return nil
	}
	out := new(XListenerSetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject satisfies runtime.Object.
func (in *XListenerSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// AddToScheme registers XListenerSet and XListenerSetList with the runtime
// scheme so the dynamic informer can decode wire JSON for the experimental
// kind. Called from pilot bootstrap.
func AddToScheme(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&XListenerSet{},
		&XListenerSetList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
