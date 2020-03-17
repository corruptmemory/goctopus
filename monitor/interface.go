package monitor

// MetricCollector sets a collector to collect metric
type MetricCollector interface {
	Name() string
	HelpMsg() string
	MetricType() uint8
	// Label returns prometheus required labels
	Label() map[string]string
	// LabelKey returns the set of keys from Label
	LabelKey() []string
	Value() float64
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
	// Register regists all the metrics on prometheus to monitor
	Register([]MetricCollector)
	// In() takes in MetricWrappers to update metrics, write only
	In() chan<- ReporterMsg
	// Start starts the reporter
	Start()
}
