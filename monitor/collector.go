package monitor

import (
	"reflect"

	"github.com/sandwich-go/boost"

	pros "github.com/sandwich-go/logbus/monitor/prometheus"

	"github.com/prometheus/client_golang/prometheus"
)

func RegisterCollector(c prometheus.Collector) {
	if c == nil {
		panic("can not register nil prometheus collector")
	}
	if isPrometheusInited() {
		boost.LogErrorAndEatError(pros.DefaultPrometheusRegistry.Register(c))
	} else {
		pros.Collectors = append(pros.Collectors, c)
	}
}

func isPrometheusInited() bool {
	m := reflect.TypeOf(DefaultMetricsReporter)
	return m.String() == "*prometheus.Reporter"
}
