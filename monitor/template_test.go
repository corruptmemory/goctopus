package monitor

import (
	"bytes"
	"testing"
)

func TestGenTemplate(t *testing.T) {
	expectedOut := `errorCountsPayload := map[string]string{
	"verb": "1",
	"endpoint": "2",
}
errorCountsMsg := reporter.GenerateMsg("errorCountsMetricName")
errorCountsMsg.SetMetricType(monitor.IncCounter)
errorCountsMsg.SetPayload(errorCountsPayload)
reporter.In() <- errorCountsMsg
`

	var buf bytes.Buffer
	testPayloadFields := []PayloadFields{
		{"verb", "1"},
		{"endpoint", "2"},
	}
	payloadName := "errorCounts"
	p, err := NewPayload("reporter", payloadName, "", testPayloadFields, IncCounter)
	if err != nil {
		t.Error("Failed to generate new payload")
	}
	err = GenPayload(p, &buf)
	if err != nil {
		t.Error("Failed to gen payload")
	}

	if buf.String() != expectedOut {
		t.Errorf("Got\n%v,\nExpected\n%v,", buf.String(), expectedOut)
	}
}

func TestGenCounterVec(t *testing.T) {
	expectedOut := `errorCountsMetric := prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "errorCountsMetricName",
		Help: "The number errors in HTTP requests.",
	},
	[]string{"verb", "endpoint", "host", "ssl", "status_code"},
)
`

	var buf bytes.Buffer
	testPayloadFields := []PayloadFields{
		{"verb", "1"},
		{"endpoint", "2"},
		{"host", "3"},
		{"ssl", "4"},
		{"status_code", "5"},
	}
	payloadName := "errorCounts"
	p, err := NewPayload("", payloadName, "The number errors in HTTP requests.", testPayloadFields, IncCounter)
	if err != nil {
		t.Error("Failed to generate new payload")
	}
	if err = GenCounterVec(p, &buf); err != nil {
		t.Errorf("Failed to generate counterVec: %v", err)
	}
	if buf.String() != expectedOut {
		t.Errorf("Got\n%v,\nExpected\n%v,", buf.String(), expectedOut)
	}
}

func TestGenMetricList(t *testing.T) {
	expectedOut := `metrics := []*prometheus.CounterVec{
	getCountsMetric,
	postCountsMetric,
	errorCountsMetric,
}
`

	var buf bytes.Buffer

	payloadName1 := "getCountsMetric"
	payloadName2 := "postCountsMetric"
	payloadName3 := "errorCountsMetric"

	if err := GenMetricList([]string{payloadName1, payloadName2, payloadName3}, &buf); err != nil {
		t.Errorf("Failed to generate counterVec: %v", err)
	}
	if buf.String() != expectedOut {
		t.Errorf("Got\n%v,\nExpected\n%v,", buf.String(), expectedOut)
	}
}

func TestGenAll(t *testing.T) {
	expectedOut := `errorCounts1Metric := prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "errorCounts1MetricName",
		Help: "The number errors in HTTP requests.",
	},
	[]string{"verb", "endpoint", "host", "ssl", "status_code"},
)
errorCounts2Metric := prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "errorCounts2MetricName",
		Help: "The number errors in HTTP requests.",
	},
	[]string{"verb", "endpoint", "host", "ssl", "status_code"},
)
errorCounts3Metric := prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "errorCounts3MetricName",
		Help: "The number errors in HTTP requests.",
	},
	[]string{"verb", "endpoint", "host", "ssl", "status_code"},
)
metrics := []*prometheus.CounterVec{
	errorCounts1Metric,
	errorCounts2Metric,
	errorCounts3Metric,
}
errorCounts1Payload := map[string]string{
	"verb": "1",
	"endpoint": "2",
	"host": "3",
	"ssl": "4",
	"status_code": "5",
}
errorCounts1Msg := reporter.GenerateMsg("errorCounts1MetricName")
errorCounts1Msg.SetMetricType(monitor.IncCounter)
errorCounts1Msg.SetPayload(errorCounts1Payload)
reporter.In() <- errorCounts1Msg
errorCounts2Payload := map[string]string{
	"verb": "1",
	"endpoint": "2",
	"host": "3",
	"ssl": "4",
	"status_code": "5",
}
errorCounts2Msg := reporter.GenerateMsg("errorCounts2MetricName")
errorCounts2Msg.SetMetricType(monitor.IncCounter)
errorCounts2Msg.SetPayload(errorCounts2Payload)
reporter.In() <- errorCounts2Msg
errorCounts3Payload := map[string]string{
	"verb": "1",
	"endpoint": "2",
	"host": "3",
	"ssl": "4",
	"status_code": "5",
}
errorCounts3Msg := reporter.GenerateMsg("errorCounts3MetricName")
errorCounts3Msg.SetMetricType(monitor.IncCounter)
errorCounts3Msg.SetPayload(errorCounts3Payload)
reporter.In() <- errorCounts3Msg
`

	var buf bytes.Buffer
	testPayloadFields := []PayloadFields{
		{"verb", "1"},
		{"endpoint", "2"},
		{"host", "3"},
		{"ssl", "4"},
		{"status_code", "5"},
	}
	payloadName := "errorCounts"
	p1, err := NewPayload("reporter", payloadName+"1", "The number errors in HTTP requests.", testPayloadFields, IncCounter)
	if err != nil {
		t.Error("Failed to generate new payload")
	}
	p2, err := NewPayload("reporter", payloadName+"2", "The number errors in HTTP requests.", testPayloadFields, IncCounter)
	if err != nil {
		t.Error("Failed to generate new payload")
	}
	p3, err := NewPayload("reporter", payloadName+"3", "The number errors in HTTP requests.", testPayloadFields, IncCounter)
	if err != nil {
		t.Error("Failed to generate new payload")
	}
	payloads := []*Payload{p1, p2, p3}
	if err = GenAll(payloads, &buf); err != nil {
		t.Errorf("Failed to generate all: %v", err)
	}
	if buf.String() != expectedOut {
		t.Errorf("Got\n%v,\nExpected\n%v,", buf.String(), expectedOut)
	}
}
