// Original source: github.com/micro/micro/blob/master/plugin/prometheus/metrics.go
package prometheus

import (
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ErrPrometheusPanic is a catch-all for the panics which can be thrown by the Prometheus client:
var ErrPrometheusPanic = errors.New("The Prometheus client panicked. Did you do something like change the tag cardinality or the type of a metric?")

// Count is a counter with key/value labels:
// New values are added to any previous one (eg "number of hits")
func (r *Reporter) Count(name string, value int64, labels prometheus.Labels) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrPrometheusPanic
		}
	}()

	counter := r.metrics.getCounter(r.stripUnsupportedCharacters(name), labels)
	metric, err := counter.GetMetricWith(r.convertLabels(labels))
	if err != nil {
		return err
	}

	metric.Add(float64(value))
	return err
}

// Gauge is a register with key/value labels:
// New values simply override any previous one (eg "current connections")
func (r *Reporter) Gauge(name string, value float64, labels prometheus.Labels) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrPrometheusPanic
		}
	}()

	gauge := r.metrics.getGauge(r.stripUnsupportedCharacters(name), labels)
	metric, err := gauge.GetMetricWith(r.convertLabels(labels))
	if err != nil {
		return err
	}

	metric.Set(value)
	return err
}

// Timing is a histogram with key/value labels:
// New values are added into a series of aggregations
func (r *Reporter) Timing(name string, value time.Duration, labels prometheus.Labels) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrPrometheusPanic
		}
	}()
	timing := r.metrics.getTiming(r.stripUnsupportedCharacters(name), labels)
	metric, err := timing.GetMetricWith(r.convertLabels(labels))
	if err != nil {
		return err
	}

	metric.Observe(value.Seconds())
	return err
}
