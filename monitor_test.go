package logbus

import (
	"bitbucket.org/funplus/sandwich/pkg/logbus/config"
	"bitbucket.org/funplus/sandwich/pkg/logbus/logreporter"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"bitbucket.org/funplus/sandwich/pkg/logbus/monitor"
	"bitbucket.org/funplus/sandwich/pkg/logbus/monitor/noop"

	. "github.com/smartystreets/goconvey/convey"
)

func TestImplementReporter(t *testing.T) {
	Convey("test noop", t, func() {
		reporter := noop.New()
		So(reporter, ShouldNotBeNil)
		So(reporter, ShouldImplement, (*monitor.Reporter)(nil))
	})
	Convey("test log reporter\n", t, func() {
		reporter := logreporter.NewLogReporter()
		So(reporter, ShouldNotBeNil)
		So(reporter, ShouldImplement, (*monitor.Reporter)(nil))
		Convey("test log monitor\n", func() {
			Init(config.NewConf(config.WithMonitorOutput(Logbus)))
			So(monitor.Count("id1", 1, map[string]string{"p": "r"}), ShouldBeNil)
			So(monitor.Gauge("id2", 2, map[string]string{"p": "q"}), ShouldBeNil)
			So(monitor.Timing("id3", time.Minute, map[string]string{"p": "s"}), ShouldBeNil)
		})
	})
	Convey("test prometheus reporter\n", t, func() {
		Init(config.NewConf(config.WithMonitorOutput(Prometheus), config.WithDefaultLabel(map[string]string{"service": "prometheus-test"})))
		labels := map[string]string{
			"tag1": "false",
			"tag2": "true",
		}

		So(monitor.Count("test.counter.1", 6, labels), ShouldBeNil)
		So(monitor.Count("test.counter.2", 19, labels), ShouldBeNil)
		So(monitor.Count("test.counter.1", 5, labels), ShouldBeNil)
		So(monitor.Gauge("test.gauge.1", 99, labels), ShouldBeNil)
		So(monitor.Gauge("test.gauge.2", 55, labels), ShouldBeNil)
		So(monitor.Gauge("test.gauge.1", 98, labels), ShouldBeNil)
		So(monitor.Timing("test.timing.1", time.Second, labels), ShouldBeNil)
		So(monitor.Timing("test.timing.2", time.Minute, labels), ShouldBeNil)

		// Test reading back the metrics:
		rsp, err := http.Get("http://localhost:9158/metrics")
		So(err, ShouldBeNil)
		So(rsp.StatusCode, ShouldEqual, http.StatusOK)

		// Read the response body and check for our metric:
		bodyBytes, err := ioutil.ReadAll(rsp.Body)
		So(err, ShouldBeNil)

		So(string(bodyBytes), ShouldContainSubstring, `test_counter_1{service="prometheus-test",tag1="false",tag2="true"} 11`)
		So(string(bodyBytes), ShouldContainSubstring, `test_counter_2{service="prometheus-test",tag1="false",tag2="true"} 19`)
		So(string(bodyBytes), ShouldContainSubstring, `test_gauge_1{service="prometheus-test",tag1="false",tag2="true"} 98`)
		So(string(bodyBytes), ShouldContainSubstring, `test_gauge_2{service="prometheus-test",tag1="false",tag2="true"} 55`)
		So(string(bodyBytes), ShouldContainSubstring, `test_timing_1{service="prometheus-test",tag1="false",tag2="true",quantile="0.5"} 1`)
		So(string(bodyBytes), ShouldContainSubstring, `test_timing_2{service="prometheus-test",tag1="false",tag2="true",quantile="0.5"} 60`)
	})
}
