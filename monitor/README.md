# Prometheus Monitoring

## Overview
This toolbox provides an API that helps to consistently report system metrics to Prometheus. 

## How to use? 

1. Initialize a reporter
2. Create a list of metrics based on event that you would like to track
    * monitor.IncCounter is the type of metric collector
    * followed by labels
3. Register the metrics
4. Start the reporter
5. Send reporter message to the reporter at where the event is triggered.
```golang
bufferSize := uint16(1000)
flushDuration := 1 * time.Second
testOut = &ReporterOut{
		make(chan ReporterMsg),
		make(chan string),
	}
// testOut could be an empty &ReporterOut{} struct
reporter := monitor.NewReporter(bufferSize, "Reporter Name", flushDuration, testOut)
// This is an example, you can use pre-defined collectors from collectorscache.go
// i.e. sampleTrackGET = GETCollector
sampleTrackGET := monitor.NewCountsCollector(
		"sample_get_counter", // name
		"The number of GET operations.", // help message
		monitor.IncCounter, // metric collector type
		[]string{"endpoint", "host"}, //labels
	)
metrics := []monitor.MetricCollector{
    sampleTrackGET,
}
reporter.Register(metrics)
go reporter.Start()
```

```golang
// Event triggers here
func GETEventTriggered(endpoint, host string, val float64){
    msg := testCollector.GenerateMsg()
    // map values to record
    counterValues := map[string]string{
		"endpoint":    endpoint,
		"host":        host,
    }
    // set the value 
    msg.SetPayload(counterValues)
    msg.SetValue(val)
    // send the reporter message to reporter
    reporter.In() <- sampleTrackGET
}
```

To catch certain event signal
```golang
// to collect ReporterMsg
<-reporter.TestOut().MsgOut
// to collect write file to disk signal 
<-reporter.TestOut().StringOut
```
