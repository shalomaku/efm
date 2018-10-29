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
	data                 efm
}

const namespace = "efm"

var efmLabels = []string{"node", "type"}

func newNodeCollector(collection efm) *nodeCollector {
	return &nodeCollector{
		vipMetric:            prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "vip_status"), "Show whether or not active vip", efmLabels, nil),
		databaseMetricStatus: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "database_status"), "Show whether or not database is UP or DOWN", efmLabels, nil),
		data:                 collection,
	}
}

func (collector *nodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.vipMetric
}

func (collector *nodeCollector) Collect(ch chan<- prometheus.Metric) {
	collector.CollectVipMetric(ch)
	collector.CollectDatabaseMetric(ch)
}

func (collector *nodeCollector) CollectVipMetric(ch chan<- prometheus.Metric) {
	var metricValue float64

	for node, detail := range collector.data.Nodes {
		if detail.VipActive {
			metricValue = 1
		} else {
			metricValue = 0
		}

		ch <- prometheus.MustNewConstMetric(collector.vipMetric, prometheus.GaugeValue, metricValue, node, detail.Type)
	}
}

func (collector *nodeCollector) CollectDatabaseMetric(ch chan<- prometheus.Metric) {
	var metricValue float64

	for node, detail := range collector.data.Nodes {
		if detail.DB == "UP" {
			metricValue = 1
		} else {
			metricValue = 0
		}

		ch <- prometheus.MustNewConstMetric(collector.databaseMetricStatus, prometheus.GaugeValue, metricValue, node, detail.Type)
	}
}
