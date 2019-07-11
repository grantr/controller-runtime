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
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

var (
	// MeasureTotalRequests is a measure which counts the total number of webhook
	// requests per webhook. It has two tags: webhook refers to the webhook path
	// and success refers to the webhook result e.g. true, false.
	MeasureTotalRequests = stats.Int64(
		"sigs.kubernetes.io/controller-runtime/measures/webhook_requests_total",
		"Total number of admission requests",
		stats.UnitNone,
	)

	// MeasureRequestLatency is a measure which keeps track of the duration
	// of webhook requests. It has one tag: webhook refers to the webhook path.
	// TODO should this be milliseconds?
	MeasureRequestLatency = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/webhook_latency_seconds",
		"Latency of processing admission requests",
		"s",
	)

	// Tag keys must conform to the restrictions described in
	// go.opencensus.io/tag/validate.go. Currently those restrictions are:
	// - length between 1 and 255 inclusive
	// - characters are printable US-ASCII

	// TagWebhook is a tag referring to the webhook path that handled a request.
	TagWebhook = mustNewTagKey("webhook")

	// TagSucceeded is a tag referring to the result of a webhook request.
	TagSucceeded = mustNewTagKey("succeeded")
)

func mustNewTagKey(k string) tag.Key {
	tagKey, err := tag.NewKey(k)
	if err != nil {
		panic(err)
	}
	return tagKey
}
