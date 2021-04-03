// Package amqplib is for encapsulating github.com/assembla/cony any operations
//
// As a quick start for publisher:
//	logger, err := loglib.NewProductionLogger()
//	if err != nil {
//		panic(err)
//	}
//	logger.Info("test logger", zap.String("hello", "world"))
package loglib

import "go.uber.org/zap"

type Logger struct {
	*zap.Logger
}

func NewProductionLogger(options ...zap.Option) (*Logger, error) {
	logger, err := zap.NewProduction(options...)
	if err != nil {
		return nil, err
	}
	return &Logger{logger}, nil
}

func NewDevelopmentLogger(options ...zap.Option) (*Logger, error) {
	logger, err := zap.NewDevelopment(options...)
	if err != nil {
		return nil, err
	}
	return &Logger{logger}, nil
}

func NewNopLogger() *Logger {
	return &Logger{zap.NewNop()}
}
