package circuit

import (
	"sync"
)

func GetCircuitManager() *CircuitManager {
	if _curcuitManager != nil {
		return _curcuitManager
	}
	_curcuitManager = &CircuitManager{
		cache: make(map[string]*CircuitConf, 0),
	}
	return _curcuitManager
}

type CircuitManager struct {
	sync.RWMutex
	cache map[CircuitID]*CircuitConf
}

func (s *CircuitManager) RegistCircuit(conf *CircuitConf) error {

	s.RLock()
	if _, ok := s.cache[conf.ID]; ok {
		s.RUnlock()
		return _errConfExisted
	}
	s.RUnlock()

	if conf.ID == "" || conf.MaxQPS <= 0 || conf.SecondsWindow <= 0 || conf.MaxFailRate <= 0 || conf.MaxFailRate >= 1 || conf.Timeout <= 0 {
		return _errConfBad
	}

	s.Lock()
	defer s.Unlock()

	s.cache[conf.ID] = conf
	return nil
}

func (s *CircuitManager) GetCircuitConf(id string) *CircuitConf {

	s.RLock()
	if conf, ok := s.cache[id]; ok {
		s.RUnlock()
		return conf
	}
	s.RUnlock()

	s.Lock()
	defer s.Unlock()

	conf := s.makeDefaultConf(id)
	s.cache[id] = conf

	return conf
}

func (s *CircuitManager) makeDefaultConf(id string) *CircuitConf {
	return &CircuitConf{
		ID:            id,
		SecondsWindow: _defaultSecondsWindow,
		MaxFailRate:   _defaultFailRate,
		Timeout:       _defaultTimeout,
		MaxQPS:        _defaultMaxQPS,
	}
}

func (s *CircuitManager) Prepare(id string, run runFn, fallback fallbackFn) *CircuitHandler {
	conf := s.GetCircuitConf(id)
	return &CircuitHandler{
		Conf:     conf,
		Run:      run,
		FallBack: fallback,
	}
}
