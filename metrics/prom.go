package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewCounter(ns, subsystem, name string, labelsOpt map[string]string) prometheus.Counter {

	if len(labelsOpt) == 0 { // no empty maps
		labelsOpt = nil
	}

	m := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
		})
	prometheus.MustRegister(m)
	return m
}

func NewCounterVec(ns, subsystem, name string, labelsOpt map[string]string, variableLabels []string) *prometheus.CounterVec {

	if len(labelsOpt) == 0 { // no empty maps
		labelsOpt = nil
	}

	m := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
		}, variableLabels)

	prometheus.MustRegister(m)
	return m
}

func NewGauge(ns, subsystem, name string, labelsOpt map[string]string) prometheus.Gauge {

	if len(labelsOpt) == 0 { // no empty maps
		labelsOpt = nil
	}

	m := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
		})
	prometheus.MustRegister(m)
	return m
}

func NewGaugeVec(ns, subsystem, name string, labelsOpt map[string]string, variableLabels []string) *prometheus.GaugeVec {
	if len(labelsOpt) == 0 { // no empty maps
		labelsOpt = nil
	}

	m := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
		}, variableLabels)
	prometheus.MustRegister(m)
	return m
}

func NewSummary(ns, subsystem, name string, labelsOpt map[string]string) prometheus.Summary {

	if len(labelsOpt) == 0 { // no empty maps
		labelsOpt = nil
	}

	m := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
		})
	prometheus.MustRegister(m)
	return m
}

func NewHistogram(ns, subsystem, name string, labelsOpt map[string]string, buckets []float64) prometheus.Histogram {

	if len(labelsOpt) == 0 { // no empty maps
		labelsOpt = nil
	}

	m := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
			Buckets:     buckets,
		})
	prometheus.MustRegister(m)
	return m
}

func NewHistogramVec(ns, subsystem, name string, labelsOpt map[string]string, buckets []float64, variableLabels []string) *prometheus.HistogramVec {
	if len(labelsOpt) == 0 {
		labelsOpt = nil
	}

	m := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
			Buckets:     buckets,
		}, variableLabels)
	prometheus.MustRegister(m)
	return m
}
