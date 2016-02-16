package metrics

import (
	"github.com/rcrowley/go-metrics"
	"net/http"
)

type metricsHandler struct {
	registry metrics.Registry
}

var Default = Handler(metrics.DefaultRegistry)

// ServeHTTP returns status 200 and writes metrics from the registry as JSON
func (h *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	metrics.WriteJSONOnce(h.registry, w)
}

// Handler creates an http.Handler that returns metrics registry r serialized as JSON
func Handler(r metrics.Registry) http.Handler {
	return &metricsHandler{registry: r}
}
