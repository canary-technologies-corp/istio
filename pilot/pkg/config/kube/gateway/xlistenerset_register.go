// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

// XListenerSet (experimental gateway.networking.x-k8s.io/v1alpha1) compat
// shim. The type was removed from sigs.k8s.io/gateway-api when it graduated
// to GA as gateway.networking.k8s.io/v1.ListenerSet in v1.5. This file
// registers a custom typed-informer for XListenerSet that bridges the
// dynamic client to the compat Go types in the xlistenersetcompat package,
// so pilot can read existing XListenerSet CRs during the 1.30 migration
// without depending on the upstream typed client (which no longer exists).
//
// Lives in the gateway package (not in xlistenersetcompat) to avoid an
// import cycle through pkg/config/schema/kubeclient.
// canary-1.30-xls patch.
package gateway

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"istio.io/istio/pilot/pkg/config/kube/gateway/xlistenersetcompat"
	"istio.io/istio/pkg/kube/krt"
)

// xlsToListenerSetCollection adapts a krt collection of experimental
// XListenerSet objects into a krt collection of the GA ListenerSet type.
// The wire shape is identical; we just rebuild the typed struct with the
// GA TypeMeta so downstream code paths (which key on TypeMeta.Kind) see a
// uniform "ListenerSet" kind. The parentKey.Kind is then rewritten back
// to gvk.XListenerSet after ListenerSetCollection runs so HTTPRoutes that
// parentRef the experimental kind still resolve their parents.
func xlsToListenerSetCollection(
	xls krt.Collection[*xlistenersetcompat.XListenerSet],
	opts krt.OptionsBuilder,
) krt.Collection[*gatewayv1.ListenerSet] {
	return krt.NewCollection(xls, func(ctx krt.HandlerContext, in *xlistenersetcompat.XListenerSet) **gatewayv1.ListenerSet {
		if in == nil {
			return nil
		}
		out := &gatewayv1.ListenerSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ListenerSet",
				APIVersion: gatewayv1.GroupVersion.String(),
			},
			ObjectMeta: *in.ObjectMeta.DeepCopy(),
			Spec:       *in.Spec.DeepCopy(),
			Status:     *in.Status.DeepCopy(),
		}
		return &out
	}, opts.WithName("XListenerSetAdapter")...)
}
