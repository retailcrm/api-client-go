package retailcrm

// BasicLogger provides basic functionality for logging.
type BasicLogger interface {
	Printf(string, ...interface{})
}

// DebugLogger can be used to easily wrap any logger with Debugf method into the BasicLogger instance.
type DebugLogger interface {
	Debugf(string, ...interface{})
}

type debugLoggerAdapter struct {
	logger DebugLogger
}

// DebugLoggerAdapter returns BasicLogger that calls underlying DebugLogger.Debugf.
func DebugLoggerAdapter(logger DebugLogger) BasicLogger {
	return &debugLoggerAdapter{logger: logger}
}

// Printf data in the log using DebugLogger.Debugf.
func (l *debugLoggerAdapter) Printf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}
