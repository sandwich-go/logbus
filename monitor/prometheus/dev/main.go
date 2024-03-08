package main

import (
	"math/rand"
	"time"

	"github.com/sandwich-go/logbus/monitor/prometheus"
)

func main() {
	reporter, err := prometheus.New(":8990",
		"/metrics",
		[]float64{0, 0.5, 0.75, 0.90, 0.95, 0.99, 1},
		map[string]string{"project": "sample"},
		time.Minute)
	if err != nil {
		panic(err)
	}
	for j := 0; j < 100; j++ {
		go func() {
			for i := 0; i < 100; i++ {
				_ = reporter.Count("test_http_count", rand.Int63n(5), map[string]string{"method": "metrics"})
				_ = reporter.Gauge("test_http_gauge", rand.Float64(), map[string]string{"method": "metrics"})
				_ = reporter.Timing("test_http_timing", time.Microsecond*time.Duration(rand.Intn(1000)), map[string]string{"method": "metrics"})
				time.Sleep(time.Second)
			}
		}()
	}
	for i := 0; i < 100; i++ {
		_ = reporter.Count("test_http_count", rand.Int63n(5), map[string]string{"method": "metrics"})
		_ = reporter.Gauge("test_http_gauge", rand.Float64(), map[string]string{"method": "metrics"})
		_ = reporter.Timing("test_http_timing", time.Microsecond*time.Duration(rand.Intn(1000)), map[string]string{"method": "metrics"})
		time.Sleep(time.Second)
	}
	wait := make(chan int, 0)
	<-wait
}
