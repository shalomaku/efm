package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type nodeInfo struct {
	Type      string `json:"type"`
	Agent     string `json:"agent"`
	DB        string `json:"db"`
	Vip       string `json:"vip"`
	VipActive bool   `json:"vip_active"`
	XLog      string `json:"xlog"`
	XLogInfo  string `json:"xloginfo"`
}

type efm struct {
	Nodes                 map[string]nodeInfo `json:"nodes"`
	AllowedNodes          []string            `json:"allowednodes"`
	Membershipcoordinator string              `json:"membershipcoordinator"`
	FailoverPriority      []string            `json:"failoverpriority"`
	MinimumStandbys       int                 `json:"minimumstandbys"`
	MissingNodes          []string            `json:"missingnodes"`
	Messages              []string            `json:"messages"`
}

type nodeCollector struct {
	vipMetric            *prometheus.Desc
	databaseMetricStatus *prometheus.Desc
	agentMetricStatus    *prometheus.Desc
	data                 efm
}

type GetMetric func(data efm) float64

const namespace = "efm"

var efmLabels = []string{"node", "type"}

func RegisterGauge(ch chan<- prometheus.Metric, desc *prometheus.Desc, data efm, getMetric func(nodeDetail nodeInfo) float64) {
	for node, detail := range data.Nodes {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, getMetric(detail), node, detail.Type)
	}
}

func newNodeCollector(collection efm) *nodeCollector {
	return &nodeCollector{
		vipMetric:            prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "vip_status"), "Show whether or not active vip", efmLabels, nil),
		databaseMetricStatus: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "database_status"), "Show whether or not database is UP or DOWN", efmLabels, nil),
		agentMetricStatus:    prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "agent_status"), "Show whether or not agent is UP or DOWN", efmLabels, nil),
		data:                 collection,
	}
}

func (collector *nodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.vipMetric
}

func (collector *nodeCollector) Collect(ch chan<- prometheus.Metric) {
	var data = collector.data

	RegisterGauge(ch, collector.vipMetric, data, func(nodeDetail nodeInfo) float64 {
		if nodeDetail.VipActive {
			return 1
		} else {
			return 0
		}
	})

	RegisterGauge(ch, collector.databaseMetricStatus, data, func(nodeDetail nodeInfo) float64 {
		if nodeDetail.DB == "UP" {
			return 1
		} else {
			return 0
		}
	})

	RegisterGauge(ch, collector.agentMetricStatus, data, func(nodeDetail nodeInfo) float64 {
		if nodeDetail.Agent == "UP" {
			return 1
		} else {
			return 0
		}
	})
}
