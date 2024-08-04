package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "number",
	}, []string{"code", "method", "func"},
)

func Init() {
	prometheus.MustRegister(httpRequestTotal)
	http.Handle("/metrics", promhttp.Handler())
	httpRequestTotal.WithLabelValues().Inc()
}
