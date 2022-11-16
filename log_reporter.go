// Original source: github.com/micro/micro/v3/metrics/logging/reporter.go
package logbus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// logReporter is an implementation of monitor.logReporter:
type logReporter struct {
	lg *StdLogger
}

// newLogReporter returns a configured logging reporter:
func newLogReporter() *logReporter {
	logger := &logReporter{
		lg: gMonitorLogger,
	}

	logger.lg.SetZLogger(logger.lg.WithOptions(zap.AddCallerSkip(2)))
	logger.lg.DebugWithChannel(Monitor, zap.String("start", "monitor will be logged"))
	return logger
}

// Count implements the monitor.logReporter interface Count method:
func (r *logReporter) Count(metricName string, value int64, labels prometheus.Labels) error {
	r.lg.InfoWithChannel(Monitor, zap.String("type", "Count"), zap.Int64(metricName, value), zap.Any("labels", labels))
	return nil
}

// Gauge implements the monitor.logReporter interface Gauge method:
func (r *logReporter) Gauge(metricName string, value float64, labels prometheus.Labels) error {
	r.lg.InfoWithChannel(Monitor, zap.String("type", "Gauge"), zap.Float64(metricName, value), zap.Any("labels", labels))
	return nil
}

// Timing implements the monitor.logReporter interface Timing method:
func (r *logReporter) Timing(metricName string, value time.Duration, labels prometheus.Labels) error {
	r.lg.InfoWithChannel(Monitor, zap.String("type", "Timing"), zap.String(metricName, value.String()), zap.Any("labels", labels))
	return nil
}
