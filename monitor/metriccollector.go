package monitor

import "github.com/prometheus/client_golang/prometheus"

// an implementation of MetricCollector
type metricCollector struct {
	name          string
	metricType    uint8
	promCollector *prometheus.CounterVec
}

func (m *metricCollector) Name() string {
	return m.name
}

// MetricType return a specific type
// For counter, there are 2 types, IncCounter or AddCounter
func (m *metricCollector) MetricType() uint8 {
	return m.metricType
}

func (m *metricCollector) Collector() *prometheus.CounterVec {
	return m.promCollector
}

func (m *metricCollector) GenerateMsg() ReporterMsg {
	return NewReporterMsg(m.name, m.metricType)
}

func NewMetricCollector(name string, metricType uint8, collector *prometheus.CounterVec) MetricCollector {

	return &metricCollector{
		name:          name,
		metricType:    metricType,
		promCollector: collector,
	}
}
