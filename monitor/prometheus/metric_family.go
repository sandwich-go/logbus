// Original source: github.com/micro/micro/blob/master/plugin/prometheus/metric_family.go
package prometheus

import (
	"fmt"
	"github.com/sandwich-go/boost/singleflight"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// metricFamily stores our cached metrics:
type metricFamily struct {
	countersSingle *singleflight.Group
	counters       sync.Map
	gaugesSingle   *singleflight.Group
	gauges         sync.Map
	timingsSingle  *singleflight.Group
	timings        sync.Map

	defaultLabels      prometheus.Labels
	prometheusRegistry *prometheus.Registry
	timingObjectives   map[float64]float64
	timingMaxAge       time.Duration
}

// newMetricFamily returns a new metricFamily (useful in case we want to change the structure later):
func (r *Reporter) newMetricFamily(defaultPercentiles []float64, defaultLabel prometheus.Labels, timingMaxAge time.Duration) metricFamily {

	// Take quantile thresholds from our pre-defined list:
	timingObjectives := make(map[float64]float64)
	for _, percentile := range defaultPercentiles {
		if quantileThreshold, ok := QuantileThresholds[percentile]; ok {
			timingObjectives[percentile] = quantileThreshold
		} else {
			panic(fmt.Sprintf("percentile %.2f not defined in prometheus: [0.0, 0.5, 0.75, 0.90, 0.95, 0.99, 1]", percentile))
		}
	}

	return metricFamily{
		countersSingle:     singleflight.New(),
		gaugesSingle:       singleflight.New(),
		timingsSingle:      singleflight.New(),
		defaultLabels:      r.convertLabels(defaultLabel),
		timingMaxAge:       timingMaxAge,
		prometheusRegistry: r.prometheusRegistry,
		timingObjectives:   timingObjectives,
	}
}

func (mf *metricFamily) listTagKeys(labels prometheus.Labels) (labelKeys []string) {
	labelKeys = make([]string, 0, len(labels))
	for key := range labels {
		labelKeys = append(labelKeys, key)
	}
	return
}

// getCounter either gets a counter, or makes a new one:
func (mf *metricFamily) getCounter(name string, labels prometheus.Labels) *prometheus.CounterVec {
	// See if we already have this counter:
	if c, ok := mf.counters.Load(name); ok {
		return c.(*prometheus.CounterVec)
	}
	c, _ := mf.countersSingle.Do(name, func() (interface{}, error) {
		// Make a new counter:
		counter := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        name,
				ConstLabels: mf.defaultLabels,
			},
			mf.listTagKeys(labels),
		)

		// Register it and add it to our list:
		mf.prometheusRegistry.MustRegister(counter)
		mf.counters.Store(name, counter)
		return counter, nil
	})
	return c.(*prometheus.CounterVec)
}

// getGauge either gets a gauge, or makes a new one:
func (mf *metricFamily) getGauge(name string, labels prometheus.Labels) *prometheus.GaugeVec {
	// See if we already have this gauge:
	if g, ok := mf.gauges.Load(name); ok {
		return g.(*prometheus.GaugeVec)
	}
	g, _ := mf.gaugesSingle.Do(name, func() (interface{}, error) {
		gauge := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        name,
				ConstLabels: mf.defaultLabels,
			},
			mf.listTagKeys(labels),
		)

		// Register it and add it to our list:
		mf.prometheusRegistry.MustRegister(gauge)
		mf.gauges.Store(name, gauge)
		return gauge, nil
	})
	return g.(*prometheus.GaugeVec)
}

// getTiming either gets a timing, or makes a new one:
func (mf *metricFamily) getTiming(name string, labels prometheus.Labels) *prometheus.SummaryVec {
	// See if we already have this timing:
	if t, ok := mf.timings.Load(name); ok {
		return t.(*prometheus.SummaryVec)
	}

	t, _ := mf.timingsSingle.Do(name, func() (interface{}, error) {
		// Make a new timing:
		timing := prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        name,
				ConstLabels: mf.defaultLabels,
				Objectives:  mf.timingObjectives,
				MaxAge:      mf.timingMaxAge,
			},
			mf.listTagKeys(labels),
		)

		// Register it and add it to our list:
		mf.prometheusRegistry.MustRegister(timing)
		mf.timings.Store(name, timing)
		return timing, nil
	})
	return t.(*prometheus.SummaryVec)
}
