package monitor

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	flushInterval     = 1 * time.Second
	trackEvent        = true
	testReporter      = NewReporter(uint16(1000), "test-reporter", flushInterval, trackEvent)
	collectorName     = "test_collector"
	mockLabel1        = "endpoint"
	mockLabel2        = "host"
	testPromCollector = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: collectorName,
			Help: "A sample test count",
		},
		[]string{mockLabel1, mockLabel2},
	)
	testMetricCollector = NewMetricCollector(collectorName, AddCounter, testPromCollector)
	msgEventChan        = testReporter.MsgEvent()
	toDiskEventChan     = testReporter.ToDiskEvent()
)

func TestIn(t *testing.T) {
	randNum := rand.New(rand.NewSource(99)).Float64()
	msg := testMetricCollector.GenerateMsg()
	// map values to record
	testValues := map[string]string{
		mockLabel1: "test-endpoint",
		mockLabel2: "test-host",
	}
	msg.SetPayload(testValues)
	msg.SetValue(randNum)
	testReporter.In() <- msg

	if msgIn, ok := <-msgEventChan; !ok {
		t.Error("Passing ReporterMsg failed")
	} else {
		if ok := reflect.DeepEqual(msgIn.Payload(), testValues); !ok {
			t.Errorf("Got wrong payload, expected %v, got %v\n", testValues, msgIn.Payload())
		}

		if msgIn.Value() != randNum {
			t.Errorf("Got wrong value, expected %v, got %v\n", randNum, msgIn.Value())
		}

		if msgIn.MetricType() != AddCounter {
			t.Errorf("Got wrong metric type, expected %v, got %v\n", IncCounter, msgIn.MetricType())
		}
	}
}

func TestWriteToFile(t *testing.T) {
	go testReporter.Start()
	randNum := rand.New(rand.NewSource(99)).Float64()
	msg := testMetricCollector.GenerateMsg()
	// map values to record
	testValues := map[string]string{
		mockLabel1: "test-endpoint",
		mockLabel2: "test-host",
	}
	msg.SetPayload(testValues)
	msg.SetValue(randNum)

	testReporter.In() <- msg
	<-msgEventChan
	stringOutMsg := <-toDiskEventChan
	if stringOutMsg != WroteFileToDiskMsg {
		t.Errorf("Got wrong msg out, expected %v, got %v\n", WroteFileToDiskMsg, stringOutMsg)
	}

	outFilePath := fmt.Sprintf("%s/%s%s", PrometheusExportDir, testReporter.Name(), PrometheusSuffix)
	if _, err := os.Stat(outFilePath); err != nil {
		t.Errorf("File %v does not exist\n", outFilePath)
	}
	_ = os.Remove(outFilePath)
}

func TestMain(m *testing.M) {
	testReporter.Register([]MetricCollector{testMetricCollector})
	code := m.Run()
	testReporter.Close()
	os.Exit(code)
}
