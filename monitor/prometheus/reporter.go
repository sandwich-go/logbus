// Original source: github.com/micro/micro/blob/master/plugin/prometheus/reporter.go
package prometheus

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/sandwich-go/boost"

	"github.com/sandwich-go/logbus/monitor/prometheus/node"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// QuantileThresholds maps quantiles / percentiles to error thresholds (required by the Prometheus client).
	// Must be from our pre-defined set [0.0, 0.5, 0.75, 0.90, 0.95, 0.99, 1]:
	QuantileThresholds = map[float64]float64{0.0: 0, 0.5: 0.05, 0.75: 0.04, 0.90: 0.03, 0.95: 0.02, 0.99: 0.001, 1: 0}
)

var Collectors []prometheus.Collector
var DefaultPrometheusRegistry = prometheus.NewRegistry()

// Reporter is an implementation of metrics.Reporter:
type Reporter struct {
	prometheusRegistry *prometheus.Registry
	metrics            metricFamily
	names              sync.Map
}

func init() {
	boost.LogErrorAndEatError(DefaultPrometheusRegistry.Register(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll),
	)))
	boost.LogErrorAndEatError(DefaultPrometheusRegistry.Register(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			Namespace: "goruntime",
		})),
	)
}

// New returns a configured prometheus reporter:
func New(
	defaultPrometheusListenAddress string,
	defaultPrometheusPath string,
	defaultPercentiles []float64,
	defaultLabel prometheus.Labels,
	timingMaxAge time.Duration) (*Reporter, error) {
	boost.LogErrorAndEatError(DefaultPrometheusRegistry.Register(node.NewNodeCollector(defaultLabel)))
	// Make a prometheus registry (this keeps track of any metrics we generate):
	DefaultPrometheusRegistry.MustRegister(Collectors...)
	Collectors = nil
	// Make a new Reporter:
	newReporter := &Reporter{
		prometheusRegistry: DefaultPrometheusRegistry,
	}
	// Add metrics families for each type:
	newReporter.metrics = newReporter.newMetricFamily(defaultPercentiles, defaultLabel, timingMaxAge)

	go func() {
		mux := http.NewServeMux()
		mux.Handle(defaultPrometheusPath, promhttp.HandlerFor(DefaultPrometheusRegistry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
		err := http.ListenAndServe(defaultPrometheusListenAddress, mux)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return newReporter, nil
}

// convertLabels turns labels into prometheus labels:
func (r *Reporter) convertLabels(labels prometheus.Labels) prometheus.Labels {
	for key, value := range labels {
		labels[key] = r.stripUnsupportedCharacters(value)
	}
	return labels
}

// listTagKeys returns a list of tag keys (we need to provide this to the Prometheus client):
func (r *Reporter) listTagKeys(labels prometheus.Labels) (labelKeys []string) {
	return r.metrics.listTagKeys(labels)
}

// stripUnsupportedCharacters cleans up a metrics key or value:
func (r *Reporter) stripUnsupportedCharacters(metricName string) string {
	name, ok := r.names.Load(metricName)
	if !ok {
		valueWithoutDots := strings.Replace(metricName, ".", "_", -1)
		valueWithoutCommas := strings.Replace(valueWithoutDots, ",", "_", -1)
		name = strings.Replace(valueWithoutCommas, " ", "", -1)
		r.names.Store(metricName, name)
	}
	return name.(string)
}
