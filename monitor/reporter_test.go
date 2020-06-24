package monitor

import (
	"bufio"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	flushInterval = 1 * time.Second
	trackEvent    = true
	testReporter  = NewReporter(uint16(1000), "test-reporter", flushInterval, trackEvent)
	collectorName = "test_collector"
	mockLabel1    = "endpoint"
	mockLabel2    = "host"

	msgEventChan    = testReporter.MsgEvent()
	toDiskEventChan = testReporter.ToDiskEvent()

	expectedSum = 0
)

func TestIn(t *testing.T) {

	N := 1000
	// map values to record
	testValues := map[string]string{
		mockLabel1: "test-endpoint",
		mockLabel2: "test-host",
	}
	expectedValues := make(map[int]int)
	for i := 0; i < N; i++ {
		randNum := rand.Intn(100)
		expectedValues[i] = randNum
		expectedSum += randNum
	}
	go func() {
		for i := 0; i < N; i++ {
			msg := testReporter.GenerateMsg(collectorName)
			msg.SetMetricType(AddCounter)
			msg.SetPayload(testValues)
			msg.SetValue(float64(expectedValues[i]))
			testReporter.In() <- msg
		}
	}()

	counter := 0
	for msgIn := range msgEventChan {

		expectedVal := expectedValues[counter]
		if ok := reflect.DeepEqual(msgIn.Payload(), testValues); !ok {
			t.Errorf("Got wrong payload, expected %v, got %v\n", testValues, msgIn.Payload())
		}
		if msgIn.Value() != float64(expectedVal) {
			t.Errorf("Got wrong value, expected %v, got %v\n", expectedVal, msgIn.Value())
		}

		if msgIn.MetricType() != AddCounter {
			t.Errorf("Got wrong metric type, expected %v, got %v\n", IncCounter, msgIn.MetricType())
		}
		counter += 1
		if counter == N-1 {
			break
		}
	}

}

func TestWriteToFile(t *testing.T) {

	N := 1000
	// map values to record
	testValues := map[string]string{
		mockLabel1: "test-endpoint",
		mockLabel2: "test-host",
	}

	for i := 0; i < N; i++ {
		randNum := rand.Intn(1000)
		msg := testReporter.GenerateMsg(collectorName)
		msg.SetMetricType(AddCounter)
		msg.SetPayload(testValues)
		msg.SetValue(float64(randNum))
		expectedSum += randNum
		testReporter.In() <- msg
	}

	<-toDiskEventChan
	counter := 0

	for range msgEventChan {
		counter += 1
		if counter == N {
			break
		}
	}

	<-toDiskEventChan

	outFilePath := fmt.Sprintf("%s/%s%s", PrometheusExportDir, testReporter.Name(), PrometheusSuffix)
	if _, err := os.Stat(outFilePath); err != nil {
		t.Errorf("File %v does not exist\n", outFilePath)
	}
	if f, err := os.Open(outFilePath); err == nil {
		reader := bufio.NewReader(f)
		b := make([]byte, 10000)
		reader.Read(b)
		lastLine := strings.SplitAfter(string(b), "\n")[2]
		result := strings.Trim(strings.SplitAfter(lastLine, " ")[1], "\n")
		flt, _, err := big.ParseFloat(result, 10, 0, big.ToNearestEven)
		if err != nil {
			panic(err)
		}
		var gotSum = new(big.Int)
		gotSum, _ = flt.Int(gotSum)

		if gotSum.Cmp(big.NewInt(int64(expectedSum))) != 0 {
			t.Errorf("Wrong result, expected %v, got %v, result file: \n %v", expectedSum, gotSum, lastLine)
		}
	}

	_ = os.Remove(outFilePath)
}

func TestMain(m *testing.M) {
	testPromCollector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: collectorName,
			Help: "A sample test count",
		},
		[]string{mockLabel1, mockLabel2},
	)
	testReporter.Register([]*prometheus.CounterVec{testPromCollector})
	go testReporter.Start()
	code := m.Run()
	testReporter.Close()
	os.Exit(code)
}
