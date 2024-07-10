package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/upstars-global/cloudflare-domain-expiration-exporter/internal/checker"
)

const namespace = "domain_expiration_checker"

type Exporter struct {
	desc    *prometheus.Desc
	checker checker.Checker
}

func New(checker checker.Checker) *Exporter {
	return &Exporter{
		checker: checker,
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "result"),
			"domain expiration check results (per domain).",
			[]string{"domain", "status"}, nil,
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.desc
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	results := e.checker.GetExpirations()
	for domain, check := range results {
		ch <- prometheus.MustNewConstMetric(e.desc, prometheus.GaugeValue, float64(check.ExpiresIn), domain, string(check.Status))
	}
}
