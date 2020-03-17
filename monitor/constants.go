package monitor

const (
	// List of prometheus metric collectors types

	// IncCounter increaments by 1
	IncCounter = uint8(0)
	// AddCoutner increaments by an arbitrary value
	AddCounter = uint8(1)
	// add from here...

	// Common collector Name
	MetricGETCountsName    = "get_counter"
	MetricGETCountsHelpMsg = "The number of GET operations."

	MetricPOSTCountsName    = "post_counter"
	MetricPOSTCountsHelpMsg = "The number of POST operations."

	MetricPixelCountsName    = "number_of_pixels_sent"
	MetricPixelCountsHelpMsg = "The total number of pixels sent"

	MetricPixelBytesName    = "bytes_of_pixels_sent"
	MetricPixelBytesHelpMsg = "The total bytes of pixels sent to stats-cache"

	MetricErrorCountName    = "number_of_errors"
	MetricErrorCountHelpMsg = "Error rate"

	MetricKafkaProduceCountName    = "number_of_msg_produced_total"
	MetricKafkaProduceCountHelpMsg = "Kafka produce rate"
	// Other constants
	PrometheusExportDir string = "/tmp/node_exporter/textfile_collector"
	PrometheusSuffix    string = ".prom"
	WroteFileToDiskMsg  string = "Wrote file to disk"
)
