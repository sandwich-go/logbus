package monitor

/*func SetDefaultMetricsReporter(
	monitorOutput config.MonitorOutput,
	defaultPrometheusListenAddress string,
	defaultPrometheusPath string,
	defaultPercentiles []float64,
	defaultLabel prometheusClient.Labels,
	timingMaxAge time.Duration) {
	switch monitorOutput {
	case config.Noop:
		DefaultMetricsReporter = noop.New()
	case config.Logbus:
		DefaultMetricsReporter = logreporter.NewLogReporter()
	case config.Prometheus:
		var err error
		DefaultMetricsReporter, err = prometheus.New(defaultPrometheusListenAddress, defaultPrometheusPath, defaultPercentiles, defaultLabel, timingMaxAge)
		if err != nil {
			panic(err)
		}
		//logbus.DebugWithChannel(logbus.Monitor, zap.String("prometheus [http] listening on", defaultPrometheusListenAddress), zap.String("path", defaultPrometheusPath))
	}
}*/
