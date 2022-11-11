package monitor

import (
	"reflect"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	pros "github.com/sandwich-go/logbus/monitor/prometheus"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterCollector(t *testing.T) {
	Convey("test reporter name", t, func() {
		m := reflect.TypeOf(DefaultMetricsReporter)
		So(m.String(), ShouldEqual, "*noop.Reporter")
		var err error
		DefaultMetricsReporter, err = pros.New(":", "/metrics", []float64{}, nil, time.Minute)
		So(err, ShouldBeNil)
		m = reflect.TypeOf(DefaultMetricsReporter)
		So(m.String(), ShouldEqual, "*prometheus.Reporter")
		So(func() { RegisterCollector(nil) }, ShouldPanic)
		RegisterCollector(prometheus.NewCounterVec(prometheus.CounterOpts{}, nil))
	})
}
