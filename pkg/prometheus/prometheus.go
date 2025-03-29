package prometheus

import (
	"slices"

	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus"
)

// HTTP
const (
	HTTPRequestsTotal      = "http_requests_total"
	HTTPRequestsDuration   = "http_requests_duration"
	HTTPRequestsSizeBytes  = "http_requests_size_bytes"
	HTTPResponsesSizeBytes = "http_responses_size_bytes"

	HelpHTTPRequestTotal       = "Total http requests amount"
	HelpHTTPRequestDuration    = "Duration of http requests"
	HelpHTTPRequestsSizeBytes  = "Tracks the size of HTTP requests"
	HelpHTTPResponsesSizeBytes = "Tracks the size of HTTP responses"
)

type Prometheus struct {
	collectors []prometheus.Collector
}

func NewPrometheus(collectors ...prometheus.Collector) *Prometheus {
	p := &Prometheus{}

	for _, collector := range collectors {
		p.addCollector(collector)
	}

	return p
}

func (p *Prometheus) NewCounter(name, help string, labels []string) *prometheus.CounterVec {
	c := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: help,
	}, labels)

	p.addCollector(c)

	return c
}

func (p *Prometheus) NewHistogram(name, help string, buckets []float64, labels []string) *prometheus.HistogramVec {
	c := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: buckets,
	}, labels)

	p.addCollector(c)

	return c
}

func (p *Prometheus) NewSummary(name, help string, labels []string) *prometheus.SummaryVec {
	c := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: name,
		Help: help,
	}, labels)

	p.addCollector(c)

	return c
}

func (p *Prometheus) NewGauge(name, help string, labels []string) *prometheus.GaugeVec {
	c := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}, labels)

	p.addCollector(c)

	return c
}

func (p *Prometheus) AddCollector(collector prometheus.Collector) {
	p.addCollector(collector)
}

func (p *Prometheus) addCollector(collector prometheus.Collector) {
	if !slices.Contains(p.collectors, collector) {
		p.collectors = append(p.collectors, collector)
	}
}

func (p *Prometheus) Register() (err error) {
	for _, collector := range p.collectors {
		err = prometheus.Register(collector)
		if err != nil {
			log.Errorf("error register collector: %v", err)
			return
		}
	}

	return nil
}
