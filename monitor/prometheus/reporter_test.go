// Original source: github.com/micro/micro/blob/master/plugin/prometheus/reporter_test.go
package prometheus

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPrometheusReporter(t *testing.T) {
	Convey("make a reporter", t, func() {
		reporter, err := New(":9999", "/prometheus",
			[]float64{0, 0.5, 0.75, 0.90, 0.95, 0.99, 1}, map[string]string{"service": "prometheus-test"}, time.Minute)
		So(err, ShouldBeNil)
		// test conversion

		labels := map[string]string{
			"tag1": "false",
			"tag2": "true",
		}
		convertedTags := reporter.convertLabels(labels)
		So(convertedTags["tag1"], ShouldEqual, "false")
		So(convertedTags["tag2"], ShouldEqual, "true")
		listedTags := reporter.listTagKeys(labels)
		So("tag1", ShouldBeIn, listedTags)
		So("tag2", ShouldBeIn, listedTags)

		// test string cleaning
		preparedMetricName := reporter.stripUnsupportedCharacters("some.kind,of tag")
		So(preparedMetricName, ShouldEqual, "some_kind_oftag")

		// test MetricFamilies
		metricFamily := reporter.newMetricFamily([]float64{0, 0.5, 0.75, 0.90, 0.95, 0.99, 1}, nil, time.Minute)
		cnt := metricFamily.getCounter("testCounter", map[string]string{"test": "test", "counter": "test"})
		So(cnt, ShouldNotBeNil)
		So(len(metricFamily.counters), ShouldEqual, 1)
		gauges := metricFamily.getGauge("testGauge", map[string]string{"test": "test", "gauge": "test"})
		So(gauges, ShouldNotBeNil)
		So(len(metricFamily.gauges), ShouldEqual, 1)
		timing := metricFamily.getTiming("testTiming", map[string]string{"test": "test", "timing": "test"})
		So(timing, ShouldNotBeNil)
		So(len(metricFamily.timings), ShouldEqual, 1)
	})
}
