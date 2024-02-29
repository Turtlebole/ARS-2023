package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// The Prometheus metrics that will be exposed.
	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total",
			Help: "Total number of http hits.",
		},
	)
	configGetHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_config_get_hit_total",
			Help: "Total number of http hits to /config/ with GET method.",
		},
	)
	configPostHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_config_post_hit_total",
			Help: "Total number of http hits to /config/ with POST method.",
		},
	)
	configDeleteHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_config_delete_hit_total",
			Help: "Total number of http hits to /config/ with DELETE method.",
		},
	)
	configsGetHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_configs_get_hit_total",
			Help: "Total number of http hits to /configs/ with GET method.",
		},
	)
	groupGetHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_group_get_hit_total",
			Help: "Total number of http hits to /group/ with GET method.",
		},
	)

	groupPostHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_group_post_hit_total",
			Help: "Total number of http hits to /group/ with POST method.",
		},
	)
	groupDeleteHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_group_delete_hit_total",
			Help: "Total number of http hits to /group/ with DELETE method.",
		},
	)
	groupPutHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_group_put_hit_total",
			Help: "Total number of http hits to /group/ with PUT method.",
		},
	)
	groupsGetHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_groups_get_hit_total",
			Help: "Total number of http hits to /groups/ with GET method.",
		},
	)

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{
		httpHits,
		configGetHits,
		configPostHits,
		configDeleteHits,
		configsGetHits,
		groupGetHits,
		groupsGetHits,
		groupPostHits,
		groupDeleteHits,
		groupPutHits,
	}

	// Prometheus Registry to register metrics.
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Register metrics that will be exposed.
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func count(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request URL path matches "/config/" and the HTTP method is "GET"
		if r.URL.Path == "/configs/" && r.Method == "GET" {
			configsGetHits.Inc()
		} else if r.URL.Path == "/config/" && r.Method == "GET" || r.Method == "POST" || r.Method == "DELETE" {
			switch r.Method {
			case "GET":
				configGetHits.Inc()
				break
			case "POST":
				configPostHits.Inc()
				break
			case "DELETE":
				configDeleteHits.Inc()
				break
			}
		} else if r.URL.Path == "/groups/" && r.Method == "GET" {
			groupsGetHits.Inc()
		} else if r.URL.Path == "/group/" && r.Method == "GET" || r.Method == "POST" || r.Method == "DELETE" || r.Method == "PUT" {
			switch r.Method {
			case "GET":
				groupGetHits.Inc()
				break
			case "POST":
				groupPostHits.Inc()
				break
			case "DELETE":
				groupDeleteHits.Inc()
				break
			case "PUT":
				groupPutHits.Inc()
				break
			}
		}
		httpHits.Inc()
	}
}
