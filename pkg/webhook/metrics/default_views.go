/*
Copyright 2018 The Kubernetes Authors.

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

package metrics

import (
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// DefaultPrometheusDistribution is an OpenCensus Distribution with the same
	// buckets as the default buckets in the Prometheus client.
	DefaultPrometheusDistribution = view.Distribution(.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10)

	// ViewTotalRequests counts TotalRequests with Webhook and Succeeded tags.
	ViewTotalRequests = view.View{
		Name:        "controller_runtime_webhook_requests_total",
		Description: "Total number of admission requests",
		Measure:     MeasureTotalRequests,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagWebhook, TagSucceeded},
	}

	// ViewRequestLatency is a histogram of RequestLatency with a Webhook tag.
	ViewRequestLatency = view.View{
		Name:        "controller_runtime_webhook_latency_seconds",
		Description: "Latency of processing admission requests",
		Measure:     MeasureRequestLatency,
		Aggregation: DefaultPrometheusDistribution,
		TagKeys:     []tag.Key{TagWebhook},
	}

	// DefaultViews is an array of OpenCensus views that can be registered
	// using view.Register(metrics.DefaultViews...) to export default metrics.
	DefaultViews = []*view.View{
		&ViewTotalRequests,
		&ViewRequestLatency,
	}
)
