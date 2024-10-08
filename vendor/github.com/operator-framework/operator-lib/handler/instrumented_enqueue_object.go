// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"

	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/operator-framework/operator-lib/handler/internal/metrics"
)

// InstrumentedEnqueueRequestForObject wraps controller-runtime handler for
// "EnqueueRequestForObject", and sets up primary resource metrics on event
// handlers. The main objective of this handler is to set prometheus metrics
// when create/update/delete events occur. These metrics contain the following
// information on resource.
//
//	resource_created_at_seconds{"name", "namespace", "group", "version", "kind"}
//
// To call the handler use:
//
//	&handler.InstrumentedEnqueueRequestForObject{}
type InstrumentedEnqueueRequestForObject[T client.Object] struct {
	handler.TypedEnqueueRequestForObject[T]
}

// Create implements EventHandler, and creates the metrics.
func (h InstrumentedEnqueueRequestForObject[T]) Create(ctx context.Context, e event.TypedCreateEvent[T], q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	setResourceMetric(e.Object)
	h.TypedEnqueueRequestForObject.Create(ctx, e, q)
}

// Update implements EventHandler, and updates the metrics.
func (h InstrumentedEnqueueRequestForObject[T]) Update(ctx context.Context, e event.TypedUpdateEvent[T], q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	setResourceMetric(e.ObjectOld)
	setResourceMetric(e.ObjectNew)

	h.TypedEnqueueRequestForObject.Update(ctx, e, q)
}

// Delete implements EventHandler, and deletes metrics.
func (h InstrumentedEnqueueRequestForObject[T]) Delete(ctx context.Context, e event.TypedDeleteEvent[T], q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	deleteResourceMetric(e.Object)
	h.TypedEnqueueRequestForObject.Delete(ctx, e, q)
}

func setResourceMetric(obj client.Object) {
	if obj != nil {
		labels := getResourceLabels(obj)
		m, _ := metrics.ResourceCreatedAt.GetMetricWith(labels)
		m.Set(float64(obj.GetCreationTimestamp().UTC().Unix()))
	}
}

func deleteResourceMetric(obj client.Object) {
	if obj != nil {
		labels := getResourceLabels(obj)
		_ = metrics.ResourceCreatedAt.Delete(labels)
	}
}

func getResourceLabels(obj client.Object) map[string]string {
	return map[string]string{
		"name":      obj.GetName(),
		"namespace": obj.GetNamespace(),
		"group":     obj.GetObjectKind().GroupVersionKind().Group,
		"version":   obj.GetObjectKind().GroupVersionKind().Version,
		"kind":      obj.GetObjectKind().GroupVersionKind().Kind,
	}
}
