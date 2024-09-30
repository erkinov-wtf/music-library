package logger

import "fmt"

func InfoLogOp(op string, msg string) {
	Logger.Info(fmt.Sprintf("%v: INF: %v", op, msg))
}

func ErrorLogOp(op string, msg string) {
	Logger.Error(fmt.Sprintf("%v: ERR: %v", op, msg))
}

func DebugLogOp(op string, msg string) {
	Logger.Debug(fmt.Sprintf("%v: DBG: %v", op, msg))
}
