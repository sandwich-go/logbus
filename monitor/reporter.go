// Original source: github.com/micro/micro/v3/metrics/go

package monitor

import (
	"time"

	"github.com/sandwich-go/logbus/monitor/noop"

	"github.com/prometheus/client_golang/prometheus"
)

// logReporter is an interface for collecting and instrumenting metrics
type Reporter interface {
	Count(metric string, value int64, labels prometheus.Labels) error
	Gauge(metric string, value float64, labels prometheus.Labels) error
	Timing(metric string, value time.Duration, labels prometheus.Labels) error
}

var (
	// DefaultMetricsReporter implementation
	DefaultMetricsReporter Reporter
)

func init() {
	DefaultMetricsReporter = noop.New()
}

// Count submits a counter metric using the DefaultMetricsReporter:
func Count(metric string, value int64, labels prometheus.Labels) error {
	return DefaultMetricsReporter.Count(metric, value, labels)
}

// Gauge submits a gauge metric using the DefaultMetricsReporter:
func Gauge(metric string, value float64, labels prometheus.Labels) error {
	return DefaultMetricsReporter.Gauge(metric, value, labels)
}

// Timing submits a timing metric using the DefaultMetricsReporter:
func Timing(metric string, value time.Duration, labels prometheus.Labels) error {
	return DefaultMetricsReporter.Timing(metric, value, labels)
}
