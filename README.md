#### Goya-Circuit: 类似于Hystrix的熔断器实现

Goya circuit is a circuit breaker mechanism implementation in microservices. It can prevent the whole link avalanche caused by the blocking of a step in the call chain. Similar implementations such as hystrix.

## Quick Start

```go
go get github.com/Anderson-Lu/goya-circuit
```

## Usage

First, set a configuration for your business like this:

```go
func registConf() {
    conf := &circuit.CircuitConf{
		ID:             "c1",
		SecondsWindow:  1,
		MaxFailRate:    0.1,
		Timeout:        1000,
		MaxQPS:         10001,
		QPSLimitOption: circuit.QPSLimitOptionFastFail,
	}
	circuit.GetCircuitManager().RegistCircuit(conf)
}
```

Then define your execution method and fallback method like this:

```go
var (
    runFun = func() error {
		return errors.New("bad")
	}
    fallbackFunc = func(error) error {		
		return nil
	}
)
```

Finally, you can safely execute via `CircuitManager` like this:

```go
func run() {
    handler := circuit.GetCircuitManager().Prepare("c1", runFun, fallbackFunc)
    handler.Go()
}
```

Or directly use external exposure methods:

```go
func run() {
    circuit.Go("c1", runFun, fallbackFunc)
}
```

## Synchronous & asynchronous

Goya circuit provides two execution methods, `Go ()` and `Do ()`. The `Go ()` method will directly start a goroutine to execute, while the `Do()` method will wait synchronously.

```go
func Do(id string, run runFn, fallback fallbackFn) {
	handler := GetCircuitManager().Prepare(id, run, fallback)
	handler.Do()
}

func Go(id string, run runFn, fallback fallbackFn) {
	go Do(id, run, fallback)
}
```

## Configuration Items

Goya circuit supports the following configuration items：

| Item           | Usage                                                        |
| -------------- | ------------------------------------------------------------ |
| ID             | ID                                                           |
| SecondsWindow  | The window size used to calculate the failure rate           |
| MaxFailRate    | Conditions for triggering fusing: failure rate               |
| Timeout        | Timeout setting for method execution                         |
| MaxQPS         | QPS current limiting                                         |
| QPSLimitOption | The operation mode after QPS current limiting supports blocking and fast failure |

## UT Coverage

```shell
Running tool: /usr/data/go1.17/go/bin/go test -timeout 30s -coverprofile=/tmp/vscode-goMPFCsu/go-code-cover goya/goya-circuit

ok  	goya/goya-circuit	12.001s	coverage: 98.8% of statements
```

## Benchmark (goya-circuit vs hystrix)

```shell
Running tool: /usr/data/go1.17/go/bin/go test -benchmem -run=^$ -coverprofile=/tmp/vscode-goMPFCsu/go-code-cover -bench . goya/goya-circuit

goos: linux
goarch: amd64
pkg: goya/goya-circuit
cpu: Intel(R) Core(TM) i5-8500 CPU @ 3.00GHz
BenchmarkHystrix-4       	  433087	      3147 ns/op	    1212 B/op	      23 allocs/op
BenchmarkGoyaCircuit-4   	 1000000	      1514 ns/op	     571 B/op	       7 allocs/op
PASS
coverage: 70.6% of statements
ok  	goya/goya-circuit	3.066s
```

