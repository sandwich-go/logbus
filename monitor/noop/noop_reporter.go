// Original source: github.com/micro/micro/v3/metrics/noop/reporter.go

package noop

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Reporter is an implementation of monitor.Reporter:
type Reporter struct {
}

// New returns a configured noop reporter:
func New() *Reporter {
	return &Reporter{}
}

// Count implements the monitor.Reporter interface Count method:
func (r *Reporter) Count(metricName string, value int64, labels prometheus.Labels) error {
	return nil
}

// Gauge implements the monitor.Reporter interface Gauge method:
func (r *Reporter) Gauge(metricName string, value float64, labels prometheus.Labels) error {
	return nil
}

// Timing implements the monitor.Reporter interface Timing method:
func (r *Reporter) Timing(metricName string, value time.Duration, labels prometheus.Labels) error {
	return nil
}
