package circuit

import (
	"context"
	"time"
)

type CircuitHandler struct {
	Conf     *CircuitConf
	Run      runFn
	FallBack fallbackFn
}

func (s *CircuitHandler) Go() {
	go s.Do()
}

func (s *CircuitHandler) Do() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(s.Conf.Timeout))
	defer cancel()

	metrics := GetMetricsManager()
	metrics.RecordDo(s.Conf)

	if needFallback, err := metrics.NeedFallback(s.Conf); needFallback {
		s.FallBack(err)
		return
	}

	done := make(chan struct{}, 0)
	fail := make(chan struct{}, 0)
	go func() {
		err := s.Run()
		if err != nil {
			fail <- struct{}{}
			return
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
		return
	case <-fail:
		metrics.RecordFail(s.Conf)
	case <-ctx.Done():
		metrics.RecordFail(s.Conf)
	}
}
