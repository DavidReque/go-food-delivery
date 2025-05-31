package models

type LogType int32

// LogType representa el tipo de logger que se utilizará
const (
	Zap    LogType = 0 // Zap es el tipo de logger que se utilizará
	Logrus LogType = 1 // Logrus es el tipo de logger que se utilizará
)
