package circuit

import (
	"testing"
	"time"
)

// func BenchmarkHystrix(b *testing.B) {
// 	hystrix.ConfigureCommand("benchmark", hystrix.CommandConfig{
// 		Timeout:               1000,
// 		MaxConcurrentRequests: 100,
// 		ErrorPercentThreshold: 50,
// 	})

// 	run := func() error {
// 		time.Sleep(time.Millisecond * 100)
// 		return nil
// 	}

// 	fallbackFn := func(error) error {
// 		return nil
// 	}

// 	b.ResetTimer()
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		hystrix.Go("benchmark", run, fallbackFn)
// 	}
// 	b.StopTimer()

// }

func BenchmarkGoyaCircuit(b *testing.B) {

	conf := &CircuitConf{
		ID:             "benchmark",
		SecondsWindow:  10,
		MaxFailRate:    0.5,
		Timeout:        1000,
		MaxQPS:         100,
		QPSLimitOption: QPSLimitOptionFastFail,
	}

	run := func() error {
		time.Sleep(time.Millisecond * 100)
		return nil
	}

	GetCircuitManager().RegistCircuit(conf)
	handler := GetCircuitManager().Prepare("benchmark", run, func(error) error { return nil })

	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		handler.Go()
	}
	b.StopTimer()
}
