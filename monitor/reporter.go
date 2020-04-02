package monitor

import (
	"encoding/json"
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

// collectorDesc to unmarshal prometheus.Desc string
type collectorDesc struct {
	Name string `json:"fqName"`
	Help string `json:"help"`
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
	if trackEvent {
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
	if r.event.StringOut != nil {
		defer close(r.event.StringOut)
	}
	writeToDisk := func() {
		r.writeToFile()
		if r.event.StringOut != nil {
			r.event.StringOut <- WroteFileToDiskMsg
		}
	}
	for {
		select {
		case <-r.done:
			writeToDisk()
			return
		case <-ticker.C:
			writeToDisk()
		}
	}
}

func (r *reporter) Register(cs []*prometheus.CounterVec) {
	descChan := make(chan *prometheus.Desc, len(cs))
	defer close(descChan)
	cd := &collectorDesc{}
	for _, c := range cs {
		r.registry.MustRegister(c)
		c.Describe(descChan)
		descriptor := <-descChan
		err := json.Unmarshal([]byte(JsonUnmarshalHelper(descriptor.String()[4:])), cd)
		if err != nil {
			log.Fatalf("Unmarshal collector description -%v- failed", c)
			panic(err)
		}
		r.counterMap[cd.Name] = c
	}
}

func (r *reporter) GenerateMsg(collectorName string) ReporterMsg {
	return NewReporterMsg(collectorName)
}

func (r *reporter) run() {
	incr := func(msg ReporterMsg) {
		if collector, ok := r.counterMap[msg.Name()]; !ok {
			log.Printf("unregistered collector %v\n", msg.Name())
		} else {
			if msg.MetricType() == AddCounter {
				collector.With(msg.Payload()).Add(msg.Value())
				return
			}
			collector.With(msg.Payload()).Inc()
		}
	}

	if r.event.MsgOut != nil {
		defer close(r.event.MsgOut)
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

func (r *reporter) Close() {
	close(r.metricsIn)
	r.done <- struct{}{}
	close(r.done)
}
