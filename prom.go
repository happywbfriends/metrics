package metrics

import "github.com/prometheus/client_golang/prometheus"

func newCounter(ns, subsystem, name string, labelsOpt map[string]string) prometheus.Counter {

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

func newCounterVec(ns, subsystem, name string, labelsOpt map[string]string, variableLabels []string) *prometheus.CounterVec {

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

func newGauge(ns, subsystem, name string, labelsOpt map[string]string) prometheus.Gauge {

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

func newSummary(ns, subsystem, name string, labelsOpt map[string]string) prometheus.Summary {

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

func newHistogram(ns, subsystem, name string, labelsOpt map[string]string, buckets []float64) prometheus.Histogram {

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
