package circuit

const (
	_defaultTimeout       = 1000
	_defaultSecondsWindow = 10
	_defaultFailRate      = 0.50
	_defaultMaxQPS        = 1 << 13

	QPSLimitOptionFastFail = 0
	QPSLimitOptionBlock    = 1
)

var (
	_curcuitManager *CircuitManager
)

type CircuitID = string
type QPSLimitOption int

type CircuitConf struct {
	ID             string
	SecondsWindow  int
	MaxFailRate    float32
	Timeout        int
	MaxQPS         int
	QPSLimitOption int
}

func NewCircuitConf(opts ...CircuitConfOption) *CircuitConf {
	conf := &CircuitConf{
		ID:            "",
		SecondsWindow: _defaultSecondsWindow,
		MaxFailRate:   _defaultFailRate,
		Timeout:       _defaultTimeout,
		MaxQPS:        _defaultMaxQPS,
	}
	for _, v := range opts {
		v(conf)
	}
	return conf
}