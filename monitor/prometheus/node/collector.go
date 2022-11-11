package node

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/net"
)

const (
	namespace         = "gonode"
	netdevSub         = "netdev"
	sockStatSubsystem = "sockstat"
)

func NewNodeCollector(defaultLabel prometheus.Labels) prometheus.Collector {
	netdevLabels := []string{"device"}
	return &nodeCollector{
		BytesSent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "bytes_sent"),
			"number of bytes sent",
			netdevLabels,
			defaultLabel,
		),
		BytesRecv: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "bytes_recv"),
			"number of bytes received",
			netdevLabels,
			defaultLabel,
		),
		PacketsSent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "packets_sent"),
			"number of packets sent",
			netdevLabels,
			defaultLabel,
		),
		PacketsRecv: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "packets_recv"),
			"number of packets received",
			netdevLabels,
			defaultLabel,
		),
		Errin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "err_in"),
			"total number of errors while receiving",
			netdevLabels,
			defaultLabel,
		),
		Errout: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "err_out"),
			"total number of errors while sending",
			netdevLabels,
			defaultLabel,
		),
		Dropin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "drop_in"),
			"total number of incoming packets which were dropped",
			netdevLabels,
			defaultLabel,
		),
		Dropout: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "drop_opt"),
			"total number of outgoing packets which were dropped (always 0 on OSX and BSD)",
			netdevLabels,
			defaultLabel,
		),
		Fifoin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "fifo_in_err"),
			"total number of FIFO buffers errors while receiving",
			netdevLabels,
			defaultLabel,
		),
		Fifoout: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netdevSub, "fifo_out_err"),
			"total number of FIFO buffers errors while sending",
			netdevLabels,
			defaultLabel,
		),
		LoadAvg: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "load"),
			"node load avg",
			[]string{"avg"},
			defaultLabel,
		),
	}
}

type nodeCollector struct {
	BytesSent   *prometheus.Desc
	Fifoout     *prometheus.Desc
	Fifoin      *prometheus.Desc
	Dropout     *prometheus.Desc
	Dropin      *prometheus.Desc
	Errout      *prometheus.Desc
	Errin       *prometheus.Desc
	PacketsRecv *prometheus.Desc
	PacketsSent *prometheus.Desc
	BytesRecv   *prometheus.Desc
	LoadAvg     *prometheus.Desc
}

func (n *nodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.BytesSent
	ch <- n.Fifoout
	ch <- n.Fifoin
	ch <- n.Dropout
	ch <- n.Dropin
	ch <- n.Errout
	ch <- n.Errin
	ch <- n.PacketsRecv
	ch <- n.PacketsSent
	ch <- n.BytesRecv
	ch <- n.LoadAvg
}

func (n *nodeCollector) Collect(ch chan<- prometheus.Metric) {
	n.collectNetdev(ch)
	n.collectSockstat(ch)
	n.collectNetstat(ch)
}

func (n *nodeCollector) collectNetdev(ch chan<- prometheus.Metric) {
	if counters, err := net.IOCounters(true); err == nil {
		for _, c := range counters {
			ch <- prometheus.MustNewConstMetric(n.BytesSent, prometheus.CounterValue, float64(c.BytesSent), c.Name)
			ch <- prometheus.MustNewConstMetric(n.BytesRecv, prometheus.CounterValue, float64(c.BytesRecv), c.Name)
			ch <- prometheus.MustNewConstMetric(n.PacketsSent, prometheus.CounterValue, float64(c.PacketsSent), c.Name)
			ch <- prometheus.MustNewConstMetric(n.PacketsRecv, prometheus.CounterValue, float64(c.PacketsRecv), c.Name)
			ch <- prometheus.MustNewConstMetric(n.Errin, prometheus.CounterValue, float64(c.Errin), c.Name)
			ch <- prometheus.MustNewConstMetric(n.Errout, prometheus.CounterValue, float64(c.Errout), c.Name)
			ch <- prometheus.MustNewConstMetric(n.Dropin, prometheus.CounterValue, float64(c.Dropin), c.Name)
			ch <- prometheus.MustNewConstMetric(n.Dropout, prometheus.CounterValue, float64(c.Dropout), c.Name)
			ch <- prometheus.MustNewConstMetric(n.Fifoin, prometheus.CounterValue, float64(c.Fifoin), c.Name)
			ch <- prometheus.MustNewConstMetric(n.Fifoout, prometheus.CounterValue, float64(c.Fifoout), c.Name)
		}
	}
	if info, err := load.Avg(); err == nil {
		ch <- prometheus.MustNewConstMetric(n.LoadAvg, prometheus.GaugeValue, float64(info.Load1), "1m")
		ch <- prometheus.MustNewConstMetric(n.LoadAvg, prometheus.GaugeValue, float64(info.Load5), "5m")
		ch <- prometheus.MustNewConstMetric(n.LoadAvg, prometheus.GaugeValue, float64(info.Load15), "15m")
	}
}
