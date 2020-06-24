package monitor

import (
	"errors"
	"io"
	"strings"
	"text/template"
)

const (
	// Counter vec template
	conterVecTemplate = `{{.Metric}} := prometheus.NewCounterVec(
{{"\t"}}prometheus.CounterOpts{
{{"\t"}}{{"\t"}}Name: "{{.MetricName}}",
{{"\t"}}{{"\t"}}Help: "{{.HelpMessage}}",
{{"\t"}}},
{{"\t"}}[]string{{"{"}}"{{join .FieldList "\", \""}}"{{"}"}},
){{"\n"}}`

	// Payload template.
	payloadTemplate = `{{.PayloadName}} := map[string]string{{"{"}}{{block "list" .Fields}}{{""}}{{range .}}
{{"\t"}}"{{print .Key }}": "{{print .Value}}",{{end}}{{end}}
{{"}"}}
{{.PayloadMsgName}} := {{.ReporterName}}.GenerateMsg("{{.MetricName}}")
{{.PayloadMsgName}}.SetMetricType({{.MetricType}})
{{.PayloadMsgName}}.SetPayload({{.PayloadName}})
{{.ReporterName}}.In() <- {{.PayloadMsgName}}{{"\n"}}`

	// Metric List template
	metricListTemplate = `metrics := []*prometheus.CounterVec{{"{"}}{{block "list" .}}{{range .}}
{{"\t"}}{{print .}},{{end}}{{end}}
{{"}"}}{{"\n"}}`
)

var (
	InvalidCounterError = errors.New("Invalid Monitor Counter Type")
	funcs               = template.FuncMap{"join": strings.Join}
)

type PayloadFields struct {
	Key   string
	Value string
}
type Payload struct {
	PayloadName    string
	PayloadMsgName string
	ReporterName   string
	Metric         string
	MetricName     string
	MetricType     string
	HelpMessage    string
	Fields         []PayloadFields
	FieldList      []string
}

// NewPayload creates a new payload for template generation
func NewPayload(reporterName, payloadName, helpMessage string, fields []PayloadFields, t CollectorType) (*Payload, error) {
	var metricType string
	if t == IncCounter {
		metricType = "monitor.IncCounter"
	} else if t == AddCounter {
		metricType = "monitor.AddCounter"
	} else {
		return &Payload{}, InvalidCounterError
	}
	var fieldList []string
	for _, field := range fields {
		fieldList = append(fieldList, field.Key)
	}

	p := &Payload{
		PayloadName:    payloadName + "Payload",
		PayloadMsgName: payloadName + "Msg",
		ReporterName:   reporterName,
		Metric:         payloadName + "Metric",
		MetricName:     payloadName + "MetricName",
		MetricType:     metricType,
		HelpMessage:    helpMessage,
		Fields:         fields,
		FieldList:      fieldList,
	}
	return p, nil
}

// GenPayload generates payload template based on Payload
// and outputs result to io.Writer provided
func GenPayload(p *Payload, writer io.Writer) (err error) {
	template, err := template.New("payload").Parse(payloadTemplate)
	if err = template.Execute(writer, p); err != nil {
		return
	}
	return
}

// GenCounterVec generate prometheus counterVec
func GenCounterVec(p *Payload, writer io.Writer) (err error) {
	template, err := template.New("counterVec").Funcs(funcs).Parse(conterVecTemplate)
	if err = template.Execute(writer, p); err != nil {
		return
	}
	return
}

func GenMetricList(metricList []string, writer io.Writer) (err error) {
	template, err := template.New("metricList").Parse(metricListTemplate)
	if err = template.Execute(writer, metricList); err != nil {
		return
	}
	return
}

// GenAll generates all payloads at once
func GenAll(payloads []*Payload, writer io.Writer) (err error) {
	var metricList []string
	for _, p := range payloads {
		err = GenCounterVec(p, writer)
		if err != nil {
			return
		}
		metricList = append(metricList, p.Metric)
	}

	err = GenMetricList(metricList, writer)
	if err != nil {
		return
	}

	for _, p := range payloads {
		err = GenPayload(p, writer)
		if err != nil {
			return
		}
	}
	return
}
