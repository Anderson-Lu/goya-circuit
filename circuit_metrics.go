package circuit

import (
	"sync"
	"time"
)

type SecondMetric struct {
	Req  int
	Fail int
	Ts   int64
}

type MetricsManager struct {
	sync.RWMutex
	cache map[CircuitID]*Metrics
}

type Metrics struct {
	sync.RWMutex
	buckets *BucketLinklist
}

var (
	_metricsManager *MetricsManager
)

func GetMetricsManager() *MetricsManager {
	if _metricsManager == nil {
		_metricsManager = &MetricsManager{
			cache: make(map[CircuitID]*Metrics, 0),
		}
	}
	return _metricsManager
}

func (s *MetricsManager) getCache(id CircuitID, defWindow int) *Metrics {
	s.RLock()
	if metric, ok := s.cache[id]; ok {
		s.RUnlock()
		return metric
	}
	s.RUnlock()

	s.Lock()
	defer s.Unlock()
	m := &Metrics{
		buckets: NewBucketLinklist(defWindow),
	}
	s.cache[id] = m
	return s.cache[id]
}

func (s *MetricsManager) RecordDo(conf *CircuitConf) {
	m := s.getCache(conf.ID, conf.SecondsWindow)
	now := time.Now().Unix()
	n := &SecondMetric{
		Ts:  now,
		Req: 1,
	}
	m.buckets.UpsertLast(n)
}

func (s *MetricsManager) RecordFail(conf *CircuitConf) {
	m := s.getCache(conf.ID, conf.SecondsWindow)
	now := time.Now().Unix()
	n := &SecondMetric{
		Ts:   now,
		Fail: 1,
	}
	m.buckets.UpsertLast(n)
}

func (s *MetricsManager) NeedFallback(conf *CircuitConf) (bool, error) {
	m := s.getCache(conf.ID, conf.SecondsWindow)

_retry:
	if conf.MaxQPS < m.buckets.GetLastQPS(time.Now().Unix()) {
		switch conf.QPSLimitOption {
		case QPSLimitOptionBlock:
			time.Sleep(time.Millisecond)
			goto _retry
		case QPSLimitOptionFastFail:
			return true, _errQPSFallback
		}
	}

	failRate := m.buckets.CalcFailRate(conf.SecondsWindow)
	if failRate >= conf.MaxFailRate {
		return true, _errFailRateFallback
	}
	return false, nil
}
