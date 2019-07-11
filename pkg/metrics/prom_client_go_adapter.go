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
	"context"
	"net/url"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	reflectormetrics "k8s.io/client-go/tools/cache"
	clientmetrics "k8s.io/client-go/tools/metrics"
	workqueuemetrics "k8s.io/client-go/util/workqueue"
)

// this file contains setup logic to initialize the myriad of places
// that client-go registers metrics.  We copy the names and formats
// from Kubernetes so that we match the core controllers.

var (
	// client metrics

	// MeasureRequestLatency ...
	MeasureRequestLatency = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/rest_client_request_latency_seconds",
		"Request latency in seconds",
		"s",
	)

	MeasureRequestResult = stats.Int64(
		"sigs.kubernetes.io/controller-runtime/measures/rest_client_requests_total",
		"Number of HTTP requests",
		stats.UnitNone,
	)

	// reflector metrics

	// TODO(directxman12): update these to be histograms once the metrics overhaul KEP
	// PRs start landing.

	MeasureListsTotal = stats.Int64(
		"sigs.kubernetes.io/controller-runtime/measures/reflector_lists_total",
		"Total number of API lists done by the reflectors",
		stats.UnitNone,
	)

	MeasureListsDuration = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/reflector_list_duration_seconds",
		"How long an API list takes to return and decode for the reflectors",
		stats.UnitNone,
	)

	// MeasureItemsPerList ... this must be a float due to the interface of SummaryMetric
	MeasureItemsPerList = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/reflector_items_per_list",
		"How many items an API list returns to the reflectors",
		stats.UnitNone,
	)

	MeasureWatchesTotal = stats.Int64(
		"sigs.kubernetes.io/controller-runtime/measures/reflector_watches_total",
		"Total number of API watches done by the reflectors",
		stats.UnitNone,
	)

	MeasureShortWatchesTotal = stats.Int64(
		"sigs.kubernetes.io/controller-runtime/measures/reflector_short_watches_total",
		"Total number of short API watches done by the reflectors",
		stats.UnitNone,
	)

	MeasureWatchDuration = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/reflector_watch_duration_seconds",
		"How long an API watch takes to return and decode for the reflectors",
		stats.UnitNone,
	)

	// MeasureItemsPerWatch ... this must be a float due to the interface of SummaryMetric
	MeasureItemsPerWatch = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/reflector_items_per_watch",
		"How many items an API watch returns to the reflectors",
		stats.UnitNone,
	)

	// MeasureLastResourceVersion ... this must be a float due to the interface of SettableGaugeMetric
	MeasureLastResourceVersion = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/reflector_last_resource_version",
		"Last resource version seen for the reflectors",
		stats.UnitNone,
	)

	// workqueue metrics

	MeasureDepth = stats.Int64(
		"sigs.kubernetes.io/controller-runtime/measures/workqueue_depth",
		"Current depth of workqueue",
		stats.UnitNone,
	)

	MeasureAdds = stats.Int64(
		"sigs.kubernetes.io/controller-runtime/measures/workqueue_adds_total",
		"Total number of adds handled by workqueue",
		stats.UnitNone,
	)

	MeasureLatency = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/workqueue_queue_latency_seconds",
		"How long in seconds an item stays in workqueue before being requested.",
		"s",
	)

	MeasureWorkDuration = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/workqueue_work_duration_seconds",
		"How long in seconds processing an item from workqueue takes.",
		"s",
	)

	MeasureRetries = stats.Int64(
		"sigs.kubernetes.io/controller-runtime/measures/workqueue_retries_total",
		"Total number of retries handled by workqueue",
		"s",
	)

	MeasureLongestRunning = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/workqueue_longest_running_processor_microseconds",
		"How many microseconds has the longest running processor for workqueue been running.",
		"us",
	)

	MeasureUnfinishedWork = stats.Float64(
		"sigs.kubernetes.io/controller-runtime/measures/workqueue_unfinished_work_seconds",
		"How many seconds of work has done that is in progress and hasn't been observed "+
			"by work_duration. Large values indicate stuck threads. One can deduce the "+
			"number of stuck threads by observing the rate at which this increases.",
		"s",
	)

	TagVerb = mustNewTagKey("verb")

	TagURL = mustNewTagKey("url")

	TagCode = mustNewTagKey("code")

	TagMethod = mustNewTagKey("method")

	TagHost = mustNewTagKey("host")

	TagName = mustNewTagKey("name")
)

func mustNewTagKey(k string) tag.Key {
	tagKey, err := tag.NewKey(k)
	if err != nil {
		panic(err)
	}
	return tagKey
}

func init() {
	clientmetrics.Register(&latencyAdapter{metric: MeasureRequestLatency}, &resultAdapter{metric: MeasureRequestResult})
	reflectormetrics.SetReflectorMetricsProvider(reflectorMetricsProvider{})
	workqueuemetrics.SetProvider(workqueueMetricsProvider{})
}

// this section contains adapters, implementations, and other sundry organic, artisinally
// hand-crafted syntax trees required to convince client-go that it actually wants to let
// someone use its metrics.

// Client metrics adapters (method #1 for client-go metrics),
// copied (more-or-less directly) from k8s.io/kubernetes setup code
// (which isn't anywhere in an easily-importable place).

type latencyAdapter struct {
	metric *stats.Float64Measure
}

func (a *latencyAdapter) Observe(verb string, u url.URL, latency time.Duration) {
	ctx, _ := tag.New(context.Background(),
		tag.Insert(TagVerb, verb),
		tag.Insert(TagURL, u.String()),
	)
	stats.Record(ctx, a.metric.M(latency.Seconds()))
}

type resultAdapter struct {
	metric *stats.Int64Measure
}

func (a *resultAdapter) Increment(code, method, host string) {
	ctx, _ := tag.New(context.Background(),
		tag.Insert(TagCode, code),
		tag.Insert(TagMethod, method),
		tag.Insert(TagHost, host),
	)
	stats.Record(ctx, a.metric.M(1))
}

// Reflector metrics provider (method #2 for client-go metrics),
// copied (more-or-less directly) from k8s.io/kubernetes setup code
// (which isn't anywhere in an easily-importable place).

type intMetric struct {
	mutators []tag.Mutator
	measure  *stats.Int64Measure
}

func (m intMetric) Inc() {
	stats.RecordWithTags(context.Background(), m.mutators, m.measure.M(1))
}

func (m intMetric) Dec() {
	stats.RecordWithTags(context.Background(), m.mutators, m.measure.M(-1))
}

type floatMetric struct {
	mutators []tag.Mutator
	measure  *stats.Float64Measure
}

func (m floatMetric) Observe(v float64) {
	stats.RecordWithTags(context.Background(), m.mutators, m.measure.M(v))
}

func (m floatMetric) Set(v float64) {
	m.Observe(v)
}

type reflectorMetricsProvider struct{}

func (reflectorMetricsProvider) NewListsMetric(name string) reflectormetrics.CounterMetric {
	return intMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureListsTotal,
	}
}

func (reflectorMetricsProvider) NewListDurationMetric(name string) reflectormetrics.SummaryMetric {
	return floatMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureListsDuration,
	}
}

func (reflectorMetricsProvider) NewItemsInListMetric(name string) reflectormetrics.SummaryMetric {
	return floatMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureItemsPerList,
	}
}

func (reflectorMetricsProvider) NewWatchesMetric(name string) reflectormetrics.CounterMetric {
	return intMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureWatchesTotal,
	}
}

func (reflectorMetricsProvider) NewShortWatchesMetric(name string) reflectormetrics.CounterMetric {
	return intMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureShortWatchesTotal,
	}
}

func (reflectorMetricsProvider) NewWatchDurationMetric(name string) reflectormetrics.SummaryMetric {
	return floatMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureWatchDuration,
	}
}

func (reflectorMetricsProvider) NewItemsInWatchMetric(name string) reflectormetrics.SummaryMetric {
	return floatMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureItemsPerWatch,
	}
}

func (reflectorMetricsProvider) NewLastResourceVersionMetric(name string) reflectormetrics.GaugeMetric {
	return floatMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureLastResourceVersion,
	}
}

// Workqueue metrics (method #3 for client-go metrics),
// copied (more-or-less directly) from k8s.io/kubernetes setup code
// (which isn't anywhere in an easily-importable place).
// TODO(directxman12): stop "cheating" and calling histograms summaries when we pull in the latest deps

type workqueueMetricsProvider struct{}

func (workqueueMetricsProvider) NewDepthMetric(name string) workqueuemetrics.GaugeMetric {
	return intMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureDepth,
	}
}

func (workqueueMetricsProvider) NewAddsMetric(name string) workqueuemetrics.CounterMetric {
	return intMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureAdds,
	}
}

func (workqueueMetricsProvider) NewLatencyMetric(name string) workqueuemetrics.SummaryMetric {
	return floatMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureLatency,
	}
}

func (workqueueMetricsProvider) NewWorkDurationMetric(name string) workqueuemetrics.SummaryMetric {
	return floatMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureWorkDuration,
	}
}

func (workqueueMetricsProvider) NewRetriesMetric(name string) workqueuemetrics.CounterMetric {
	return intMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureRetries,
	}
}

func (workqueueMetricsProvider) NewLongestRunningProcessorMicrosecondsMetric(name string) workqueuemetrics.SettableGaugeMetric {
	return floatMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureLongestRunning,
	}
}

func (workqueueMetricsProvider) NewUnfinishedWorkSecondsMetric(name string) workqueuemetrics.SettableGaugeMetric {
	return floatMetric{
		mutators: []tag.Mutator{tag.Insert(TagName, name)},
		measure:  MeasureUnfinishedWork,
	}
}
