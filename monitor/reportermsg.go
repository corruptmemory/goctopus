package monitor

// Implement ReporterMsg

// TODO: plan to convert reporterMsg type into a RealtimeResult type.

type reporterMsg struct {
	name       string
	data       map[string]string
	value      float64
	metricType uint8
}

func NewReporterMsg(name string) ReporterMsg {
	return &reporterMsg{
		name: name,
	}
}

func (r *reporterMsg) Name() string {
	return r.name
}

func (r *reporterMsg) Payload() map[string]string {
	return r.data
}

// SetPayload set payload in a verbose way, order insensitive.
func (r *reporterMsg) SetPayload(payload map[string]string) {
	r.data = payload
}

func (r *reporterMsg) Value() float64 {
	return r.value
}

func (r *reporterMsg) SetValue(val float64) {
	r.value = val
}

func (r *reporterMsg) SetMetricType(metricTyp uint8) {
	r.metricType = metricTyp
}

func (r *reporterMsg) MetricType() uint8 {
	return r.metricType
}

// Clone returns a deep copy of self
func (r *reporterMsg) Clone() ReporterMsg {
	newMsg := reporterMsg{
		name:       r.name,
		metricType: r.metricType,
	}
	if r.Payload() != nil {
		newMsg.SetPayload(r.Payload())
	}
	if r.Value() != 0 {
		newMsg.SetValue(r.Value())
	}
	return &newMsg
}
