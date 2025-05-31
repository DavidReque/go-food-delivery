package models

type LogType int32

// LogType representa el tipo de logger que se utilizará
const (
	Zap    LogType = 0
	Logrus LogType = 1
)