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

// how frequent the monitor write to disk
flushInterval := 15 * time.Second
// if you don't want to track event signals
// set this to false
trackEvent = true
reporter := monitor.NewReporter(bufferSize, "Reporter Name", flushInterval, trackEvent)

// channels for event signals
msgEventChan        = reporter.MsgEvent()
toDiskEventChan     = reporter.ToDiskEvent()

// Create a prometheus CounterVec
CollectorName := "test_counter"
samplePromCollector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: CollectorName,
			Help: "A sample test count",
		},
		[]string{"label1", "label2"},
    )

// Register metrics
metrics := []*prometheus.CounterVec{
    samplePromCollector,
    // more metrics ...
}
reporter.Register(metrics)

// start the reporter
go reporter.Start()
```

```golang
// Event triggers here
func GETEventTriggered(endpoint, host string, val float64){
    // generate a reporterMsg based on metric name
    msg := reporter.GenerateMsg(CollectorName)
    // map values to record
    counterValues := map[string]string{
		"label1":    endpoint,
		"label2":        host,
    }
    // set a metric type
    // AddCounter or IncCounter
    msg.SetMetricType(monitor.AddCounter)
    // set the value 
    msg.SetPayload(counterValues)
    // if using a monitor.AddCounter, don't forget to call SetValue()
    msg.SetValue(val)
    // send the reporter message to reporter
    // if don't send message to reporter, the message won't be tracked
    reporter.In() <- msg
}
```

To catch reporter event signals
```golang
go func(){
    // a channel of ReporterMsg
    for sendMsgEvent := range msgEventChan {
        // do something here...
    }
}()

go func(){
    // a channel of String
    for writeToDiskEvent := range toDiskEventChan {
        // do something here...
    }
}()

```
