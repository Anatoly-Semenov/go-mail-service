package logger

import "sync"

var (
	globalLogger *Logger
	once         sync.Once
)

func InitGlobalLogger(level string, isProduction bool) error {
	var err error
	once.Do(func() {
		globalLogger, err = NewLogger(level, isProduction)
	})
	return err
}

func GetLogger() *Logger {
	if globalLogger == nil {
		
		_ = InitGlobalLogger("info", false)
	}
	return globalLogger
}
