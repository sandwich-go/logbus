//go:build !nonetstat
// +build !nonetstat

// Original source: https://github.com/prometheus/node_exporter
package node

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

var (
	netStatsSubsystem = "netstat"
	// Used for calculating the total memory bytes on TCP and UDP.
	pageSize = os.Getpagesize()
)

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

func (n *nodeCollector) collectNetstat(ch chan<- prometheus.Metric) {
	netStats, err := getNetStats("/proc/net/snmp")
	if err != nil {
		return
	}
	for protocol, protocolStats := range netStats {
		if protocol != "Tcp" && protocol != "Udp" {
			continue
		}
		for name, value := range protocolStats {
			key := protocol + "_" + name
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, netStatsSubsystem, key),
					fmt.Sprintf("Statistic %s.", protocol+name),
					nil, nil,
				),
				prometheus.UntypedValue, v,
			)
		}
	}
}

func getNetStats(fileName string) (map[string]map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseNetStats(file, fileName)
}

func parseNetStats(r io.Reader, fileName string) (map[string]map[string]string, error) {
	var (
		netStats = map[string]map[string]string{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		nameParts := strings.Split(scanner.Text(), " ")
		scanner.Scan()
		valueParts := strings.Split(scanner.Text(), " ")
		// Remove trailing :.
		protocol := nameParts[0][:len(nameParts[0])-1]
		netStats[protocol] = map[string]string{}
		if len(nameParts) != len(valueParts) {
			return nil, fmt.Errorf("mismatch field count mismatch in %s: %s",
				fileName, protocol)
		}
		for i := 1; i < len(nameParts); i++ {
			netStats[protocol][nameParts[i]] = valueParts[i]
		}
	}

	return netStats, scanner.Err()
}
