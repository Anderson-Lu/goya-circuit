package circuit

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCircuitManager(t *testing.T) {

	// singleton mode
	m1 := GetCircuitManager()
	m2 := GetCircuitManager()
	assert.Equal(t, m1, m2)

	c1 := &CircuitConf{
		ID:            "c1",
		SecondsWindow: 10,
		Timeout:       10,
		MaxFailRate:   0.1,
		MaxQPS:        1,
	}
	err := m1.RegistCircuit(c1)
	assert.Nil(t, err)

	err = m1.RegistCircuit(c1)
	assert.NotNil(t, err)

	// invalid param
	c1.ID = ""
	err = m1.RegistCircuit(c1)
	assert.NotNil(t, err)

	// get existed conf
	x1 := m1.GetCircuitConf("c1")
	assert.Equal(t, c1, x1)

	// get non-exited conf
	x2 := m1.GetCircuitConf("c2")
	assert.NotNil(t, x2)

	var reqCount = 0
	var fallbackCount = 0

	// prepare
	handler := m1.Prepare(c1.ID, func() error {
		reqCount++
		return nil
	}, func(error) error {
		fallbackCount++
		return nil
	})

	handler.Do()
	handler.Do()
	handler.Do()

	assert.Equal(t, reqCount, 3)
	assert.Equal(t, fallbackCount, 0)

	var wg sync.WaitGroup
	// qps limit
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			handler.Do()
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 103, reqCount)
	assert.Equal(t, 0, fallbackCount)

	handler.Go()
}
