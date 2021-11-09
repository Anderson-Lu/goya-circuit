package circuit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitMetrics(t *testing.T) {
	// sigleton object
	m1 := GetMetricsManager()
	assert.NotNil(t, _metricsManager)
	m2 := GetMetricsManager()
	assert.Equal(t, m1, m2)

	// get conf doen't exists or exists
	c1 := m1.getCache("c1", 1)
	c2 := m1.getCache("c1", 1)
	assert.Equal(t, c1, c2)

	// conf
	conf := &CircuitConf{
		ID:            "c1",
		SecondsWindow: 10,
		MaxFailRate:   0.5,
		Timeout:       1000,
		MaxQPS:        100,
	}

	// Record qps
	m1.RecordDo(conf)
	m1.RecordDo(conf)
	m1.RecordDo(conf)
	qps := m1.getCache(conf.ID, 10).buckets.GetLastQPS(time.Now().Unix())
	assert.Equal(t, 3, qps)

	// Recodr fail
	m1.RecordFail(conf)
	needFallback, _ := m1.NeedFallback(conf)
	assert.Equal(t, false, needFallback)
	assert.Equal(t, float32(1.0/3.0), m1.getCache(conf.ID, 10).buckets.CalcFailRate(10))

	m1.RecordFail(conf)
	needFallback, _ = m1.NeedFallback(conf)
	assert.Equal(t, true, needFallback)
	assert.Equal(t, float32(2.0/3.0), m1.getCache(conf.ID, 10).buckets.CalcFailRate(10))

	m1.RecordFail(conf)
	needFallback, _ = m1.NeedFallback(conf)
	assert.Equal(t, true, needFallback)
	assert.Equal(t, float32(3.0/3.0), m1.getCache(conf.ID, 1).buckets.CalcFailRate(10))

	// conf
	conf2 := &CircuitConf{
		ID:            "c2",
		SecondsWindow: 10,
		MaxFailRate:   0.5,
		Timeout:       1000,
		MaxQPS:        1,
	}

	m1.RecordDo(conf2)
	m1.RecordDo(conf2)
	m1.RecordDo(conf2)

	needFallback, _ = m1.NeedFallback(conf2)
	assert.Equal(t, true, needFallback)
	qps2 := m1.getCache(conf2.ID, 10).buckets.GetLastQPS(time.Now().Unix())
	assert.Equal(t, 3, qps2)

	// empty
	needFallback, _ = m1.NeedFallback(&CircuitConf{
		ID:            "c3",
		SecondsWindow: 10,
		MaxFailRate:   0.5,
		Timeout:       1000,
		MaxQPS:        1,
	})
	assert.Equal(t, false, needFallback)

	conf3 := &CircuitConf{
		ID:             "c3",
		SecondsWindow:  10,
		MaxFailRate:    0.5,
		Timeout:        1,
		MaxQPS:         1,
		QPSLimitOption: QPSLimitOptionBlock,
	}
	m1.RecordDo(conf3)
	m1.RecordDo(conf3)
	m1.RecordDo(conf3)
	needFallback, _ = m1.NeedFallback(conf3)
	assert.Equal(t, false, needFallback)
}
