//go:build !linux
// +build !linux

package node

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

func (n *nodeCollector) collectNetstat(ch chan<- prometheus.Metric) {

}

// Used for calculating the total memory bytes on TCP and UDP.
var pageSize = os.Getpagesize()

func (n *nodeCollector) collectSockstat(ch chan<- prometheus.Metric) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return
	}
	s, err := fs.NetSockstat()
	if err != nil {
		return
	}
	if s.Used != nil {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, sockStatSubsystem, "sockets_used"),
				"Number of IPv4 sockets in use.",
				nil,
				nil,
			),
			prometheus.GaugeValue,
			float64(*s.Used),
		)
	}
	// A name and optional value for a sockstat metric.
	type ssPair struct {
		name string
		v    *int
	}

	// Previously these metric names were generated directly from the file output.
	// In order to keep the same level of compatibility, we must map the fields
	// to their correct names.
	for _, p := range s.Protocols {
		pairs := []ssPair{
			{
				name: "inuse",
				v:    &p.InUse,
			},
			{
				name: "orphan",
				v:    p.Orphan,
			},
			{
				name: "tw",
				v:    p.TW,
			},
			{
				name: "alloc",
				v:    p.Alloc,
			},
			{
				name: "mem",
				v:    p.Mem,
			},
			{
				name: "memory",
				v:    p.Memory,
			},
		}

		// Also export mem_bytes values for sockets which have a mem value
		// stored in pages.
		if p.Mem != nil {
			v := *p.Mem * pageSize
			pairs = append(pairs, ssPair{
				name: "mem_bytes",
				v:    &v,
			})
		}

		for _, pair := range pairs {
			if pair.v == nil {
				// This value is not set for this protocol; nothing to do.
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(
						namespace,
						sockStatSubsystem,
						fmt.Sprintf("%s_%s", p.Protocol, pair.name),
					),
					fmt.Sprintf("Number of %s sockets in state %s.", p.Protocol, pair.name),
					nil,
					nil,
				),
				prometheus.GaugeValue,
				float64(*pair.v),
			)
		}
	}
}
