// pkg/internal/logging/logging.go (100% complete)
package logging

import "log"

type Logger interface {
    Debugf(format string, args ...interface{})
    Infof(format string, args ...interface{})  
    Warnf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
    Fatalf(format string, args ...interface{})  
}

type NopLogger struct{}

func NewNopLogger() *NopLogger {
    return &NopLogger{}
}

func (l *NopLogger) Debugf(format string, args ...interface{}) {}
func (l *NopLogger) Infof(format string, args ...interface{})  {}
func (l *NopLogger) Warnf(format string, args ...interface{})  {}
func (l *NopLogger) Errorf(format string, args ...interface{}) {}
func (l *NopLogger) Fatalf(format string, args ...interface{}) {}

type StdLogger struct {
    debug *log.Logger
    info  *log.Logger
    warn  *log.Logger
    err   *log.Logger
    fatal *log.Logger  
}

func NewStdLogger(debug, info, warn, err, fatal *log.Logger) *StdLogger {
    return &StdLogger{
        debug: debug,
        info:  info,
        warn:  warn, 
        err:   err,
        fatal: fatal,
    }
}

func (l *StdLogger) Debugf(format string, args ...interface{}) {
    l.debug.Printf(format, args...)
}

func (l *StdLogger) Infof(format string, args ...interface{}) {
    l.info.Printf(format, args...)
}

func (l *StdLogger) Warnf(format string, args ...interface{}) {
    l.warn.Printf(format, args...)  
}

func (l *StdLogger) Errorf(format string, args ...interface{}) {
    l.err.Printf(format, args...)
}

func (l *StdLogger) Fatalf(format string, args ...interface{}) {
    l.fatal.Printf(format, args...)
}
