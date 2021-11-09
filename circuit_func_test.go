package circuit

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	Do("k1", func() error {
		return nil
	}, func(error) error {
		return nil
	})
}

func TestGo(t *testing.T) {
	Go("k1", func() error {
		return nil
	}, func(e error) error {
		return nil
	})
}

func TestFallback(t *testing.T) {
	k1 := &CircuitConf{
		ID:            "k2",
		SecondsWindow: 10,
		Timeout:       10000,
		MaxFailRate:   0.1,
		MaxQPS:        1,
	}
	m1 := GetCircuitManager()
	err := m1.RegistCircuit(k1)
	assert.Nil(t, err)

	var fallbackCount = 0

	run := func() error {
		return errors.New("bad")
	}

	fallbackFn := func(error) error {
		fallbackCount++
		return nil
	}

	Do("k2", run, fallbackFn)
	Do("k2", run, fallbackFn)
	Do("k2", run, fallbackFn)

	assert.Equal(t, 2, fallbackCount)
}

func TestTimeout(t *testing.T) {
	k1 := &CircuitConf{
		ID:             "k3",
		SecondsWindow:  10,
		Timeout:        1000,
		MaxFailRate:    0.01,
		MaxQPS:         10,
		QPSLimitOption: QPSLimitOptionFastFail,
	}
	m1 := GetCircuitManager()
	err := m1.RegistCircuit(k1)
	assert.Nil(t, err)

	var fallbackCount = 0

	var wg sync.WaitGroup

	run := func() error {
		time.Sleep(time.Second * 2)
		wg.Done()
		return nil
	}

	fallbackFn := func(error) error {
		wg.Done()
		fallbackCount++
		return nil
	}

	wg.Add(1)
	go Do("k2", run, fallbackFn)

	time.Sleep(time.Second)
	wg.Add(1)
	go Do("k2", run, fallbackFn)

	time.Sleep(time.Second)
	wg.Add(1)
	go Do("k2", run, fallbackFn)

	wg.Wait()

	assert.Equal(t, 3, fallbackCount)
}
