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
	flushInterval time.Duration
	done          chan struct{}
	// event is used to observe send, write to disk events
	event *ReporterEvent
}

type ReporterEvent struct {
	MsgOut    chan ReporterMsg
	StringOut chan string
}

func NewReporter(bufferSize uint16, reporterName string, flushInterval time.Duration, trackEvent bool) *reporter {
	in := make(chan ReporterMsg, bufferSize)
	c := make(map[string]*prometheus.CounterVec)
	registry := prometheus.NewRegistry()
	var event *ReporterEvent
	if trackEvent == true {
		event = &ReporterEvent{
			make(chan ReporterMsg),
			make(chan string),
		}
	} else {
		event = &ReporterEvent{}
	}

	r := reporter{
		reporterName:  reporterName,
		metricsIn:     in,
		counterMap:    c,
		registry:      registry,
		flushInterval: flushInterval,
		done:          make(chan struct{}),
		event:         event,
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
	ticker := time.NewTicker(r.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-r.done:
			return
		case <-ticker.C:
			r.writeToFile()
			if r.event.StringOut != nil {
				r.event.StringOut <- WroteFileToDiskMsg
			}
		}
	}
}

func (r *reporter) Register(metrics []MetricCollector) {
	for _, metric := range metrics {
		r.counterMap[metric.Name()] = metric.Collector()
		r.registry.MustRegister(metric.Collector())
	}
}

func (r *reporter) run() {
	incr := func(msg ReporterMsg) {
		if collector, ok := r.counterMap[msg.Name()]; !ok {
			log.Printf("unregistered collector %v\n", msg.Name())
			return
		} else {
			if msg.MetricType() == AddCounter {
				collector.With(msg.Payload()).Add(msg.Value())
				return
			}
			collector.With(msg.Payload()).Inc()
		}

	}

	for msg := range r.metricsIn {
		if r.event.MsgOut != nil {
			r.event.MsgOut <- msg.Clone()
		}
		incr(msg)
	}
}

// ToDiskEvent observes write to disk events
func (r *reporter) ToDiskEvent() <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for msg := range r.event.StringOut {
			out <- msg
		}
	}()
	return out
}

// MsgEvent observes message passing on to Prometheus
func (r *reporter) MsgEvent() <-chan ReporterMsg {
	out := make(chan ReporterMsg)
	go func() {
		defer close(out)
		for msg := range r.event.MsgOut {
			out <- msg
		}
	}()
	return out
}

func (r *reporter) DrainAndStop() {
	close(r.metricsIn)
	close(r.event.MsgOut)
	close(r.event.StringOut)
	r.done <- struct{}{}
	close(r.done)
}
