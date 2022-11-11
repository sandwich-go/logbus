package monitor

import (
	"reflect"

	"bitbucket.org/funplus/sandwich/base/slog"

	pros "bitbucket.org/funplus/sandwich/pkg/logbus/monitor/prometheus"

	"github.com/prometheus/client_golang/prometheus"
)

func RegisterCollector(c prometheus.Collector) {
	if c == nil {
		panic("can not register nil prometheus collector")
	}
	if isPrometheusInited() {
		slog.LogErrorAndEatError(pros.DefaultPrometheusRegistry.Register(c))
	} else {
		pros.Collectors = append(pros.Collectors, c)
	}
}

func isPrometheusInited() bool {
	m := reflect.TypeOf(DefaultMetricsReporter)
	return m.String() == "*prometheus.Reporter"
}
