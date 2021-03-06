package api

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/go-logr/logr"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type echoLogrLogger struct {
	logr   logr.Logger
	output io.Writer
	prefix string
	v      log.Lvl
}

func newEchoLogrLogger(logr logr.Logger) echo.Logger {
	return &echoLogrLogger{
		logr:   logr,
		output: nopWriter{},
	}
}

type nopWriter struct{}

func (w nopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (e echoLogrLogger) Output() io.Writer {
	return e.output
}

func (e *echoLogrLogger) SetOutput(w io.Writer) {
	e.output = w
}

func (e echoLogrLogger) Prefix() string {
	return e.prefix
}

func (e *echoLogrLogger) SetPrefix(p string) {
	e.prefix = p
}

func (e echoLogrLogger) Level() log.Lvl {
	return e.v
}

func (e *echoLogrLogger) SetLevel(v log.Lvl) {
	e.v = v
}

func (e echoLogrLogger) SetHeader(string) {}

func (e echoLogrLogger) Print(i ...interface{}) {
	e.Info(i...)
}

func (e echoLogrLogger) Printf(format string, args ...interface{}) {
	e.Infof(format, args...)
}

func (e echoLogrLogger) Printj(j log.JSON) {
	e.Infoj(j)
}

func (e echoLogrLogger) Debug(i ...interface{}) {
	e.logr.V(int(e.v) + 1).Info(fmt.Sprint(i...))
}

func (e echoLogrLogger) Debugf(format string, args ...interface{}) {
	e.logr.V(int(e.v) + 1).Info(fmt.Sprintf(format, args...))
}

func (e echoLogrLogger) Debugj(j log.JSON) {
	kvs := make([]interface{}, 0)
	for k, v := range j {
		kvs = append(kvs, k, v)
	}

	e.logr.V(int(e.v)+1).WithName(e.prefix).Info("", kvs...)
}

func (e echoLogrLogger) Info(i ...interface{}) {
	e.logr.V(int(e.v)).WithName(e.prefix).Info(fmt.Sprint(i...))
}

func (e echoLogrLogger) Infof(format string, args ...interface{}) {
	e.logr.V(int(e.v)).WithName(e.prefix).Info(fmt.Sprintf(format, args...))
}

func (e echoLogrLogger) Infoj(j log.JSON) {
	kvs := make([]interface{}, 0)
	for k, v := range j {
		kvs = append(kvs, k, v)
	}

	e.logr.V(int(e.v)).WithName(e.prefix).Info("", kvs...)
}

func (e echoLogrLogger) Warn(i ...interface{}) {
	e.logr.V(int(e.v)+1).WithName(e.prefix).Error(nil, fmt.Sprint(i...))
}

func (e echoLogrLogger) Warnf(format string, args ...interface{}) {
	e.logr.V(int(e.v)+1).WithName(e.prefix).Error(nil, fmt.Sprintf(format, args...))
}

func (e echoLogrLogger) Warnj(j log.JSON) {
	kvs := make([]interface{}, 0)
	for k, v := range j {
		kvs = append(kvs, k, v)
	}

	e.logr.V(int(e.v)+1).WithName(e.prefix).Error(nil, "", kvs...)
}

func (e echoLogrLogger) Error(i ...interface{}) {
	e.logr.V(int(e.v)).WithName(e.prefix).Error(nil, fmt.Sprint(i...))
}

func (e echoLogrLogger) Errorf(format string, args ...interface{}) {
	e.logr.V(int(e.v)).WithName(e.prefix).Error(nil, fmt.Sprintf(format, args...))
}

func (e echoLogrLogger) Errorj(j log.JSON) {
	kvs := make([]interface{}, 0)
	for k, v := range j {
		kvs = append(kvs, k, v)
	}

	e.logr.V(int(e.v)).WithName(e.prefix).Error(nil, "", kvs...)
}

func (e echoLogrLogger) Fatal(i ...interface{}) {
	e.Error(i...)
	os.Exit(1)
}

func (e echoLogrLogger) Fatalj(j log.JSON) {
	e.Errorj(j)
	os.Exit(1)
}

func (e echoLogrLogger) Fatalf(format string, args ...interface{}) {
	e.Errorf(format, args...)
	os.Exit(1)
}

func (e echoLogrLogger) Panic(i ...interface{}) {
	e.Error(i...)
	panic(fmt.Sprint(i...))
}

func (e echoLogrLogger) Panicj(j log.JSON) {
	e.Errorj(j)

	txt, err := json.Marshal(j)
	if err != nil {
		panic("(echoLogrLogger).Panicj called (but it failed to serialize provided JSON)")
	}
	panic(txt)
}

func (e echoLogrLogger) Panicf(format string, args ...interface{}) {
	e.Errorf(format, args...)
	panic(fmt.Sprintf(format, args...))
}
