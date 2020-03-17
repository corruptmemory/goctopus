package monitor

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	testOut = &ReporterOut{
		make(chan ReporterMsg),
		make(chan string),
	}
	flushDuration       = 1 * time.Second
	reporterWithTestOut = NewReporter(uint16(1000), "test-reporter", flushDuration, testOut)
	testCollector       = GETCollector
)

func TestIn(t *testing.T) {
	randNum := rand.New(rand.NewSource(99)).Float64()
	msg := testCollector.GenerateMsg()
	// map values to record
	testValues := map[string]string{
		"endpoint": "test-endpoint",
		"host":     "test-host",
	}
	msg.SetPayload(testValues)
	msg.SetValue(randNum)

	reporterWithTestOut.In() <- msg

	if msgIn, ok := <-reporterWithTestOut.TestOut().MsgOut; !ok {
		t.Error("Passing ReporterMsg failed")
	} else {
		if ok := reflect.DeepEqual(msgIn.Payload(), testValues); !ok {
			t.Errorf("Got wrong payload, expected %v, got %v\n", testValues, msgIn.Payload())
		}

		if msgIn.Value() != randNum {
			t.Errorf("Got wrong value, expected %v, got %v\n", randNum, msgIn.Value())
		}

		if msgIn.MetricType() != IncCounter {
			t.Errorf("Got wrong metric type, expected %v, got %v\n", IncCounter, msgIn.MetricType())
		}
	}
}

func TestWriteToFile(t *testing.T) {
	go reporterWithTestOut.Start()
	randNum := rand.New(rand.NewSource(99)).Float64()
	msg := testCollector.GenerateMsg()
	// map values to record
	testValues := map[string]string{
		"endpoint": "test-endpoint",
		"host":     "test-host",
	}
	msg.SetPayload(testValues)
	msg.SetValue(randNum)

	reporterWithTestOut.In() <- msg
	<-reporterWithTestOut.TestOut().MsgOut
	stringOutMsg := <-reporterWithTestOut.TestOut().StringOut
	if stringOutMsg != WroteFileToDiskMsg {
		t.Errorf("Got wrong msg out, expected %v, got %v\n", WroteFileToDiskMsg, stringOutMsg)
	}

	outFilePath := fmt.Sprintf("%s/%s%s", PrometheusExportDir, reporterWithTestOut.Name(), PrometheusSuffix)
	if _, err := os.Stat(outFilePath); err != nil {
		t.Errorf("File %v does not exist\n", outFilePath)
	}
	_ = os.Remove(outFilePath)
}

func TestMain(m *testing.M) {
	reporterWithTestOut.Register([]MetricCollector{testCollector})
	code := m.Run()
	reporterWithTestOut.DrainAndStop()
	os.Exit(code)
}
