package monitor

type CollectorType uint8

const (
	// List of prometheus metric collectors types

	// IncCounter increase by 1
	IncCounter CollectorType = iota
	// AddCounter increase by an arbitrary value
	AddCounter
	// add from here...

	// Other constants
	PrometheusExportDir string = "/tmp/node_exporter/textfile_collector"
	PrometheusSuffix    string = ".prom"
	WroteFileToDiskMsg  string = "Wrote file to disk"
)
