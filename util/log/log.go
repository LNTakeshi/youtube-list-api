package log

import (
	"context"
	"fmt"

	"github.com/glassonion1/logz"
)

type Logger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Warningf(ctx context.Context, format string, args ...interface{})
	Criticalf(ctx context.Context, format string, args ...interface{})
}

type logger struct {
}

type localLogger struct {
}

func NewLogger() Logger {
	return &logger{}
}

func NewlocalLogger() Logger {
	return &localLogger{}
}

func (l *logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	logz.Debugf(ctx, format, args...)
}

func (l *logger) Infof(ctx context.Context, format string, args ...interface{}) {
	logz.Infof(ctx, format, args...)
}

func (l *logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	logz.Errorf(ctx, format, args...)
}

func (l *logger) Warningf(ctx context.Context, format string, args ...interface{}) {
	logz.Warningf(ctx, format, args...)
}

func (l *logger) Criticalf(ctx context.Context, format string, args ...interface{}) {
	logz.Criticalf(ctx, format, args...)
}

func (l *localLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *localLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *localLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *localLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *localLogger) Criticalf(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
