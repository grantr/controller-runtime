package metrics

import (
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// DefaultPrometheusDistribution is an OpenCensus Distribution with the same
	// buckets as the default buckets in the Prometheus client.
	DefaultPrometheusDistribution = view.Distribution(.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10)

	// ViewRequestLatency ...
	ViewRequestLatency = view.View{
		Name:        "rest_client_request_latency_seconds",
		Description: "Request latency in seconds. Broken down by verb and URL.",
		Measure:     MeasureRequestLatency,
		// equivalent to prometheus.ExponentialBuckets(0.001, 2, 10)
		Aggregation: view.Distribution(.001, .002, 0.004, .008, .016, .032, .064, .128, .256, .512),
		TagKeys:     []tag.Key{TagVerb, TagURL},
	}

	ViewRequestResult = view.View{
		Name:        "rest_client_requests_total",
		Description: "Number of HTTP requests, partitioned by status code, method, and host.",
		Measure:     MeasureRequestResult,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagCode, TagMethod, TagHost},
	}

	ViewListsTotal = view.View{
		Name:        "reflector_lists_total",
		Description: "Total number of API lists done by the reflectors",
		Measure:     MeasureListsTotal,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagName},
	}

	ViewListsDuration = view.View{
		Name:        "reflector_list_duration_seconds",
		Description: "How long an API list takes to return and decode for the reflectors",
		Measure:     MeasureListsDuration,
		// TODO Converted from SummaryVec. Determine correct distribution.
		Aggregation: DefaultPrometheusDistribution,
		TagKeys:     []tag.Key{TagName},
	}

	ViewItemsPerList = view.View{
		Name:        "reflector_items_per_list",
		Description: "How many items an API list returns to the reflectors",
		Measure:     MeasureItemsPerList,
		// TODO Converted from SummaryVec. Determine correct distribution.
		Aggregation: view.Distribution(0, 1, 5, 10, 50, 100),
		TagKeys:     []tag.Key{TagName},
	}

	ViewWatchesTotal = view.View{
		Name:        "reflector_watches_total",
		Description: "Total number of API watches done by the reflectors",
		Measure:     MeasureWatchesTotal,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagName},
	}

	ViewShortWatchesTotal = view.View{
		Name:        "reflector_short_watches_total",
		Description: "Total number of short API watches done by the reflectors",
		Measure:     MeasureShortWatchesTotal,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagName},
	}

	ViewWatchesDuration = view.View{
		Name:        "reflector_watch_duration_seconds",
		Description: "How long an API watch takes to return and decode for the reflectors",
		Measure:     MeasureWatchDuration,
		// TODO Converted from SummaryVec. Determine correct distribution.
		Aggregation: DefaultPrometheusDistribution,
		TagKeys:     []tag.Key{TagName},
	}

	ViewItemsPerWatch = view.View{
		Name:        "reflector_items_per_watch",
		Description: "How many items an API watch returns to the reflectors",
		Measure:     MeasureItemsPerWatch,
		// TODO Converted from SummaryVec. Determine correct distribution.
		Aggregation: view.Distribution(0, 1, 5, 10, 50, 100),
		TagKeys:     []tag.Key{TagName},
	}

	ViewLastResourceVersion = view.View{
		Name:        "reflector_last_resource_version",
		Description: "Last resource version seen for the reflectors",
		Measure:     MeasureLastResourceVersion,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{TagName},
	}

	ViewDepth = view.View{
		Name:        "workqueue_depth",
		Description: "Current depth of workqueue",
		Measure:     MeasureDepth,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{TagName},
	}

	ViewAdds = view.View{
		Name:        "workqueue_adds_total",
		Description: "Total number of adds handled by workqueue",
		Measure:     MeasureAdds,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagName},
	}

	ViewLatency = view.View{
		Name:        "workqueue_queue_latency_seconds",
		Description: "How long in seconds an item stays in workqueue before being requested.",
		Measure:     MeasureLatency,
		// equivalent to prometheus.ExponentialBuckets(10e-9, 10, 10)
		Aggregation: view.Distribution(1e-08, 1e-07, 1e-06, 1e-05, 1e-04, 0.001, 0.01, 0.1, 1, 10),
		TagKeys:     []tag.Key{TagName},
	}

	ViewWorkDuration = view.View{
		Name:        "workqueue_work_duration_seconds",
		Description: "How long in seconds processing an item from workqueue takes.",
		Measure:     MeasureWorkDuration,
		// equivalent to prometheus.ExponentialBuckets(10e-9, 10, 10)
		Aggregation: view.Distribution(1e-08, 1e-07, 1e-06, 1e-05, 1e-04, 0.001, 0.01, 0.1, 1, 10),
		TagKeys:     []tag.Key{TagName},
	}

	ViewRetries = view.View{
		Name:        "workqueue_retries_total",
		Description: "Total number of retries handled by workqueue",
		Measure:     MeasureRetries,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagName},
	}

	ViewLongestRunning = view.View{
		Name:        "workqueue_longest_running_processor_microseconds",
		Description: "How many microseconds has the longest running processor for workqueue been running.",
		Measure:     MeasureLongestRunning,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{TagName},
	}

	ViewUnfinishedWork = view.View{
		Name: "workqueue_unfinished_work_seconds",
		Description: "How many seconds of work has done that is in progress and hasn't been observed " +
			"by work_duration. Large values indicate stuck threads. One can deduce the " +
			"number of stuck threads by observing the rate at which this increases.",
		Measure:     MeasureUnfinishedWork,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{TagName},
	}

	// DefaultViews is an array of OpenCensus views that can be registered
	// using view.Register(metrics.DefaultViews...) to export default metrics.
	DefaultViews = []*view.View{
		&ViewRequestLatency,
		&ViewRequestResult,
		&ViewListsTotal,
		&ViewListsDuration,
		&ViewItemsPerList,
		&ViewWatchesTotal,
		&ViewShortWatchesTotal,
		&ViewWatchesDuration,
		&ViewItemsPerWatch,
		&ViewLastResourceVersion,
		&ViewDepth,
		&ViewAdds,
		&ViewLatency,
		&ViewWorkDuration,
		&ViewRetries,
		&ViewLongestRunning,
		&ViewUnfinishedWork,
	}
)
