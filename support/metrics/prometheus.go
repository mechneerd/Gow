package metrics

import (
	"net/http"
	"strconv"
)

// PrometheusExporter provides basic /metrics endpoint.
type PrometheusExporter struct {
	metrics map[string]float64
}

func NewPrometheusExporter() *PrometheusExporter {
	return &PrometheusExporter{metrics: make(map[string]float64)}
}

func (p *PrometheusExporter) Inc(name string) {
	p.metrics[name]++
}

func (p *PrometheusExporter) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for k, v := range p.metrics {
			w.Write([]byte(k + " " + strconv.FormatFloat(v, 'f', -1, 64) + "\n"))
		}
	}
}

