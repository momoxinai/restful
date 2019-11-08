// Package base implements a common library for Abac APIs.
// Copyright 2019 Policy Center Author. All Rights Reserved.
// The license belongs to Platform Team.
// Version 1.0
package log

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"strings"

	"wwwin-github.cisco.com/DevNet/restful/config"

	log "github.com/sirupsen/logrus"
)

// ContextLogger is the name of logger in context
const ContextLogger = "ContextLogger"

// LogOption defines the cli options for service
type LogOption struct {
	LogFormat string `long:"log-format" description:"text or json" default:"json" env:"LOG_FORMAT"`
	LogLevel  string `long:"log-level" description:"debug, info, warn, error, or fatal" default:"info" env:"LOG_LEVEL"`
}

// Logger defines a generic logging interface
type Logger interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	WithField(key string, value interface{}) Logger
}

// LogOpt is global logging configuration instance
var (
	hostname, _  = os.Hostname()
	globalEntry  = log.StandardLogger().WithFields(log.Fields{"host": hostname})
	globalLogger = &logrusLogger{skip: 4, entry: globalEntry}
)

// NewLogger will create a new Logger instance for service
func NewLogger(service string) Logger {
	if strings.Compare(config.String("log_format"), "json") == 0 {
		log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05"})
	} else {
		textFmt := log.TextFormatter{
			ForceColors:      false,
			DisableColors:    true,
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  "2006-01-02T15:04:05",
			DisableSorting:   false,
		}
		log.SetFormatter(&textFmt)
	}
	log.SetOutput(os.Stderr)
	lvl, err := log.ParseLevel(config.String("log_level"))
	if err != nil {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(lvl)
	}
	host := os.Getenv("HOSTNAME")
	if host == "" {
		host = hostname
	}
	globalEntry = log.WithFields(log.Fields{
		"service": service,
		"host":    host,
	})
	globalLogger = &logrusLogger{skip: 4, entry: globalEntry}
	return logrusLogger{skip: 4, entry: globalEntry}
}

// LoggerFromRequest will return a Logger instance from request object
func LoggerFromRequest(r *http.Request) Logger {
	if r == nil {
		return globalLogger
	}
	return LoggerFromContext(r.Context())
}

// LoggerFromContext will return a Logger instance from context object
func LoggerFromContext(ctx context.Context) Logger {
	if ctx == nil {
		return globalLogger
	}
	logger := ctx.Value(ContextLogger)
	if logger == nil {
		return globalLogger
	}
	return logger.(Logger)
}

// LogErrorf logs message in error level
func LogErrorf(format string, v ...interface{}) {
	globalLogger.Errorf(format, v...)
}

// LogInfof logs message in info level
func LogInfof(format string, v ...interface{}) {
	globalLogger.Infof(format, v...)
}

// LogDebugf logs message in debug level
func LogDebugf(format string, v ...interface{}) {
	globalLogger.Debugf(format, v...)
}

// LogFatalf logs message in fatal level
func LogFatalf(format string, v ...interface{}) {
	globalLogger.Fatalf(format, v...)
}

// LogWarnf logs message in warn level
func LogWarnf(format string, v ...interface{}) {
	globalLogger.Warnf(format, v...)
}

// LogPrintf logs message in error level
func LogPrintf(format string, v ...interface{}) {
	globalLogger.Printf(format, v...)
}

type logrusLogger struct {
	skip  int
	entry *log.Entry
}

func (l logrusLogger) entryWithCallInfo() *log.Entry {
	file, line, fn := l.callerInfo()
	return l.entry.WithFields(log.Fields{"file": file, "line": line, "func": fn})
}

func (l logrusLogger) Debugf(format string, v ...interface{}) {
	l.entryWithCallInfo().Debugf(format, v...)
}

func (l logrusLogger) Fatalf(format string, v ...interface{}) {
	l.entryWithCallInfo().Fatalf(format, v...)
}

func (l logrusLogger) Infof(format string, v ...interface{}) {
	l.entryWithCallInfo().Infof(format, v...)
}

func (l logrusLogger) Errorf(format string, v ...interface{}) {
	l.entryWithCallInfo().Errorf(format, v...)
}

func (l logrusLogger) Warnf(format string, v ...interface{}) {
	l.entryWithCallInfo().Warnf(format, v...)
}

func (l logrusLogger) Printf(format string, v ...interface{}) {
	l.entryWithCallInfo().Printf(format, v...)
}

func (l logrusLogger) Print(v ...interface{}) {
	l.entryWithCallInfo().Print(v...)
}

func (l logrusLogger) Println(v ...interface{}) {
	l.entryWithCallInfo().Println(v...)
}

func (l logrusLogger) WithField(key string, value interface{}) Logger {
	return logrusLogger{
		skip:  l.skip + 1,
		entry: l.entry.WithFields(log.Fields{key: value}),
	}
}

// callerInfo Retrieve caller info, file name and line number
func (l logrusLogger) callerInfo() (file string, line int, fn string) {
	pc, file, line, _ := runtime.Caller(l.skip)
	index := strings.Index(file, "pkg")
	if index != -1 {
		file = file[index+4:] // truncate text before pkg
	}
	fn = runtime.FuncForPC(pc).Name()
	index = strings.Index(fn, "pkg")
	if index != -1 {
		fn = fn[index+4:] // truncate text before pkg
	}
	return file, line, fn
}
