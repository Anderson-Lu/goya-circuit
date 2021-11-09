package circuit

import "strconv"

type GoyaCircuitError struct {
	Code    int
	Message string
}

func (s GoyaCircuitError) Error() string {
	return "[Goya Circuit Error] (Code:" + strconv.Itoa(s.Code) + ")" + s.Message
}

var (
	_errConfExisted      = &GoyaCircuitError{Code: 900, Message: "Regist circuit failed, conf exists"}
	_errConfBad          = &GoyaCircuitError{Code: 901, Message: "Bad circuit config params found"}
	_errQPSFallback      = &GoyaCircuitError{Code: 902, Message: "QPS limited"}
	_errFailRateFallback = &GoyaCircuitError{Code: 903, Message: "Failure rate exceed definition"}
)
