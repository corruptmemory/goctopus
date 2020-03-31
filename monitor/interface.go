package monitor

import "github.com/prometheus/client_golang/prometheus"

// MetricCollector sets a collector to collect metric
type MetricCollector interface {
	Name() string
	Collector() *prometheus.CounterVec
	MetricType() uint8
	// GenerateMsg generates a ReporterMsg to pass data to reporter
	GenerateMsg() ReporterMsg
}

// ReporterMsg carries the data from event triggered to reporter
type ReporterMsg interface {
	Name() string
	Payload() map[string]string
	SetPayload(map[string]string)
	Value() float64
	SetValue(float64)
	MetricType() uint8
	Clone() ReporterMsg
}

// Reporter takes in an ReporterMsg
type Reporter interface {
	Name() string
	// Register registers all the metric collectors on prometheus
	Register([]MetricCollector)
	// In() takes in MetricWrappers to update metrics, write only
	In() chan<- ReporterMsg
	// ToDiskEvent returns a channel of string, which observes write to disk events
	ToDiskEvent() <-chan string
	// MsgEvent observes messages passing on to prometheus
	MsgEvent() <-chan ReporterMsg
	// Start starts the reporter
	Start()
	// Shutdown the reporter
	Close()
}
