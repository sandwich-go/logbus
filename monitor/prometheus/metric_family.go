// Original source: github.com/micro/micro/blob/master/plugin/prometheus/metric_family.go
package prometheus

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// metricFamily stores our cached metrics:
type metricFamily struct {
	counterMutex sync.Mutex
	counters     map[string]*prometheus.CounterVec
	gaugeMutex   sync.Mutex
	gauges       map[string]*prometheus.GaugeVec
	timingMutex  sync.Mutex
	timings      map[string]*prometheus.SummaryVec

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
		counters:           make(map[string]*prometheus.CounterVec),
		gauges:             make(map[string]*prometheus.GaugeVec),
		timings:            make(map[string]*prometheus.SummaryVec),
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
	mf.counterMutex.Lock()
	defer mf.counterMutex.Unlock()

	// See if we already have this counter:
	counter, ok := mf.counters[name]
	if !ok {

		// Make a new counter:
		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        name,
				ConstLabels: mf.defaultLabels,
			},
			mf.listTagKeys(labels),
		)

		// Register it and add it to our list:
		mf.prometheusRegistry.MustRegister(counter)
		mf.counters[name] = counter
	}

	return counter
}

// getGauge either gets a gauge, or makes a new one:
func (mf *metricFamily) getGauge(name string, labels prometheus.Labels) *prometheus.GaugeVec {
	mf.gaugeMutex.Lock()
	defer mf.gaugeMutex.Unlock()

	// See if we already have this gauge:
	gauge, ok := mf.gauges[name]
	if !ok {

		// Make a new gauge:
		gauge = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        name,
				ConstLabels: mf.defaultLabels,
			},
			mf.listTagKeys(labels),
		)

		// Register it and add it to our list:
		mf.prometheusRegistry.MustRegister(gauge)
		mf.gauges[name] = gauge
	}

	return gauge
}

// getTiming either gets a timing, or makes a new one:
func (mf *metricFamily) getTiming(name string, labels prometheus.Labels) *prometheus.SummaryVec {
	mf.timingMutex.Lock()
	defer mf.timingMutex.Unlock()

	// See if we already have this timing:
	timing, ok := mf.timings[name]
	if !ok {

		// Make a new timing:
		timing = prometheus.NewSummaryVec(
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
		mf.timings[name] = timing
	}

	return timing
}
