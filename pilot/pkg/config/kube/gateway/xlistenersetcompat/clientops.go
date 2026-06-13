// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

// List/Watch operations for the experimental XListenerSet kind, implemented
// via the dynamic client + JSON conversion to typed objects. Used by the
// generated kubeclient dispatch to avoid panicking for XListenerSet GVR after
// upstream gateway-api dropped it in v1.5.
// canary-1.30-xls patch.
package xlistenersetcompat

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

var GroupVersionResource = SchemeGroupVersion.WithResource("xlistenersets")

// List returns all XListenerSet objects in the namespace via dynamic client.
// The unstructured list is JSON-round-tripped into the typed XListenerSetList.
func List(d dynamic.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	raw, err := d.Resource(GroupVersionResource).Namespace(namespace).List(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("listing xlistenersets: %w", err)
	}
	out := &XListenerSetList{}
	out.SetResourceVersion(raw.GetResourceVersion())
	out.Items = make([]XListenerSet, 0, len(raw.Items))
	for i := range raw.Items {
		item, err := unstructuredToXLS(&raw.Items[i])
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		out.Items = append(out.Items, *item)
	}
	return out, nil
}

// Watch streams XListenerSet events via the dynamic client, converting each
// Unstructured payload into the typed XListenerSet on the fly.
func Watch(d dynamic.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	raw, err := d.Resource(GroupVersionResource).Namespace(namespace).Watch(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("watching xlistenersets: %w", err)
	}
	return newConvertingWatch(raw), nil
}

func unstructuredToXLS(u *unstructured.Unstructured) (*XListenerSet, error) {
	buf, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("marshal unstructured: %w", err)
	}
	var item XListenerSet
	if err := json.Unmarshal(buf, &item); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &item, nil
}

type convertingWatch struct {
	source watch.Interface
	out    chan watch.Event
}

func newConvertingWatch(source watch.Interface) *convertingWatch {
	w := &convertingWatch{source: source, out: make(chan watch.Event)}
	go w.run()
	return w
}

func (w *convertingWatch) Stop() { w.source.Stop() }

func (w *convertingWatch) ResultChan() <-chan watch.Event { return w.out }

func (w *convertingWatch) run() {
	defer close(w.out)
	for ev := range w.source.ResultChan() {
		u, ok := ev.Object.(*unstructured.Unstructured)
		if !ok {
			w.out <- ev
			continue
		}
		converted, err := unstructuredToXLS(u)
		if err != nil {
			w.out <- ev
			continue
		}
		ev.Object = converted
		w.out <- ev
	}
}
