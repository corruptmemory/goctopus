package monitor

import "github.com/prometheus/client_golang/prometheus"

// ReporterMsg carries the data from event triggered to reporter
type ReporterMsg interface {
	Name() string
	Payload() map[string]string
	SetPayload(map[string]string)
	Value() float64
	SetValue(float64)
	SetMetricType(CollectorType)
	MetricType() CollectorType
	Clone() ReporterMsg
}

// Reporter takes in an ReporterMsg
type Reporter interface {
	Name() string
	// Register registers all the metrics on prometheus
	Register([]*prometheus.CounterVec)
	// GenerateMsg generates a ReporterMsg to pass data to reporter
	GenerateMsg(string) ReporterMsg
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
