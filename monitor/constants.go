package monitor

const (
	// List of prometheus metric collectors types

	// IncCounter increase by 1
	IncCounter = uint8(0)
	// AddCounter increase by an arbitrary value
	AddCounter = uint8(1)
	// add from here...

	// Other constants
	PrometheusExportDir string = "/tmp/node_exporter/textfile_collector"
	PrometheusSuffix    string = ".prom"
	WroteFileToDiskMsg  string = "Wrote file to disk"
)
