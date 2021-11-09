package circuit

type runFn func() error
type fallbackFn func(error) error

func Do(id string, run runFn, fallback fallbackFn) {
	handler := GetCircuitManager().Prepare(id, run, fallback)
	handler.Do()
}

func Go(id string, run runFn, fallback fallbackFn) {
	go Do(id, run, fallback)
}
