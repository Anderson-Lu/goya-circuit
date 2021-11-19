package circuit

type CircuitConfOption func(*CircuitConf)

func WithID(id string) CircuitConfOption {
	return func(c *CircuitConf) {
		c.ID = id
	}
}

func WithSecondWindow(secondWindow int) CircuitConfOption {
	return func(c *CircuitConf) {
		c.SecondsWindow = secondWindow
	}
}

func WithMaxFailRate(maxFailRate float32) CircuitConfOption {
	return func(c *CircuitConf) {
		c.MaxFailRate = maxFailRate
	}
}

func WithTimeout(timeout int) CircuitConfOption {
	return func(c *CircuitConf) {
		c.Timeout = timeout
	}
}

func WithMaxQPS(maxPQs int) CircuitConfOption {
	return func(c *CircuitConf) {
		c.MaxQPS = maxPQs
	}
}

func WithQPSLimitOption(limitOption int) CircuitConfOption {
	return func(c *CircuitConf) {
		c.QPSLimitOption = limitOption
	}
}
