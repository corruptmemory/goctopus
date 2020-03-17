package monitor

// a set of commonly used collectors
var (
	// Track number of GET operations
	GETCollector = NewCountsCollector(
		MetricGETCountsName,          // name
		MetricGETCountsHelpMsg,       // help message
		IncCounter,                   // metric collector type
		[]string{"endpoint", "host"}, //labels
	)
	// Track number of POST operations
	POSTCollector = NewCountsCollector(
		MetricPOSTCountsName,
		MetricPOSTCountsHelpMsg,
		IncCounter,
		[]string{"endpoint", "host"},
	)
	// Track number of pixels sent
	PixelsCountCollector = NewCountsCollector(
		MetricPixelCountsName,
		MetricPixelCountsHelpMsg,
		IncCounter,
		[]string{"port", "status"},
	)
	// Tracking total size of pixels sent
	PixelsByteCollector = NewCountsCollector(
		MetricPixelBytesName,
		MetricPixelBytesHelpMsg,
		AddCounter,
		[]string{"port"},
	)
	// Track error rate
	ErrorCountCollector = NewCountsCollector(
		MetricErrorCountName,
		MetricErrorCountHelpMsg,
		IncCounter,
		[]string{"type", "detail"},
	)
	// Track Kafka produce rate
	KafkaProduceCountCollector = NewCountsCollector(
		MetricKafkaProduceCountName,
		MetricKafkaProduceCountHelpMsg,
		IncCounter,
		[]string{"status"},
	)
)
