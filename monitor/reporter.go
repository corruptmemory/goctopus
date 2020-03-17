package monitor

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type reporter struct {
	// reporterName is used as the file name of in text collector
	reporterName  string
	metricsIn     chan ReporterMsg
	counterMap    map[string]*prometheus.CounterVec
	registry      *prometheus.Registry
	flushDuration time.Duration
	done          chan struct{}
	// reporterOut is used in test only
	testOut *ReporterOut
}

type ReporterOut struct {
	MsgOut    chan ReporterMsg
	StringOut chan string
}

func NewReporter(bufferSize uint16, reporterName string, flushDuration time.Duration, testOut *ReporterOut) *reporter {
	in := make(chan ReporterMsg, bufferSize)
	c := make(map[string]*prometheus.CounterVec)
	registry := prometheus.NewRegistry()
	r := reporter{
		reporterName:  reporterName,
		metricsIn:     in,
		counterMap:    c,
		registry:      registry,
		flushDuration: flushDuration,
		done:          make(chan struct{}),
		testOut:       testOut,
	}

	go r.run()
	return &r
}

func (r *reporter) Name() string {
	return r.reporterName
}

// In take in a ReporterMsg, write only
func (r *reporter) In() chan<- ReporterMsg {
	return r.metricsIn
}

// Out is for test purpose only, read only
func (r *reporter) Out() <-chan ReporterMsg {
	return r.metricsIn
}

func (r *reporter) writeToFile() {
	// get a clean export directory
	err := EnsureDir(PrometheusExportDir, true)
	if err != nil {
		log.Println(err)
	}
	fn := fmt.Sprintf("%s/%s%s", PrometheusExportDir, r.reporterName, PrometheusSuffix)

	if _, err := os.Create(fn); err != nil {
		log.Println(err)
	}

	if err := prometheus.WriteToTextfile(fn, r.registry); err != nil {
		log.Println(err)
	}
}

func (r *reporter) Start() {
	ticker := time.NewTicker(r.flushDuration)
	defer ticker.Stop()

	for {
		select {
		case <-r.done:
			return
		case <-ticker.C:
			r.writeToFile()
			if r.testOut.StringOut != nil {
				r.testOut.StringOut <- WroteFileToDiskMsg
			}
		}
	}
}

func (r *reporter) Register(metrics []MetricCollector) {
	for _, metric := range metrics {
		metricCount := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metric.Name(),
				Help: metric.HelpMsg(),
			},
			metric.LabelKey(),
		)
		r.counterMap[metric.Name()] = metricCount
		r.registry.MustRegister(metricCount)
	}
}

func (r *reporter) run() {
	incr := func(msg ReporterMsg) {
		if msg.MetricType() == AddCounter {
			r.counterMap[msg.Name()].With(msg.Payload()).Add(msg.Value())
			return
		}
		r.counterMap[msg.Name()].With(msg.Payload()).Inc()
	}

	for msg := range r.metricsIn {
		if r.testOut.MsgOut != nil {
			r.testOut.MsgOut <- msg.Clone()
		}
		incr(msg)
	}
}

// TestOut() provide channels to send out signals
// test use only
func (r *reporter) TestOut() *ReporterOut {
	return r.testOut
}

func (r *reporter) DrainAndStop() {
	close(r.metricsIn)
	r.done <- struct{}{}
	close(r.done)
}
