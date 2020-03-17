package monitor

// CountsCollector is an implementation of MetricCollector
type CountsCollector struct {
	name       string
	helpMsg    string
	metricType uint8
	label      map[string]string
	labelKeys  []string
	value      float64
}

func (c *CountsCollector) Name() string {
	return c.name
}

func (c *CountsCollector) HelpMsg() string {
	return c.helpMsg
}

// MetricType return a specific type
// For counter, there are 2 types, IncCounter or AddCounter
func (c *CountsCollector) MetricType() uint8 {
	return c.metricType
}

func (c *CountsCollector) Label() map[string]string {
	return c.label
}

func (c *CountsCollector) LabelKey() []string {
	return c.labelKeys
}

func (c *CountsCollector) SetMapValue(k, v string) {
	c.label[k] = v
}

func (c *CountsCollector) SetMap(m map[string]string) {
	for k, v := range m {
		c.SetMapValue(k, v)
	}
}

func (c *CountsCollector) Value() float64 {
	return c.value
}

func (c *CountsCollector) SetValue(v float64) {
	c.value = v
}

func (c *CountsCollector) GenerateMsg() ReporterMsg {
	return NewReporterMsg(c.name, c.metricType)
}

func NewCountsCollector(name, helpMsg string, metricType uint8, labelKeys []string) MetricCollector {
	l := make(map[string]string)
	for _, labelKey := range labelKeys {
		l[labelKey] = ""
	}
	return &CountsCollector{
		name:       name,
		helpMsg:    helpMsg,
		metricType: metricType,
		label:      l,
		labelKeys:  labelKeys,
	}
}
