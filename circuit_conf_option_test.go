package circuit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConfWithOptions(t *testing.T) {
	conf := NewCircuitConf(
		WithID("2221"),
		WithMaxFailRate(0.3),
		WithMaxQPS(1000),
		WithQPSLimitOption(0),
		WithSecondWindow(2),
		WithTimeout(1000),
	)
	assert.Equal(t, "2221", conf.ID)
	assert.Equal(t, float32(0.3), conf.MaxFailRate)
	assert.Equal(t, 1000, conf.MaxQPS)
	assert.Equal(t, 0, conf.QPSLimitOption)
	assert.Equal(t, 2, conf.SecondsWindow)
	assert.Equal(t, 1000, conf.Timeout)
}
