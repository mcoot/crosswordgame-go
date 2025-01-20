package logging

import (
	"context"
	"github.com/mcoot/crosswordgame-go/internal/utils"
	"go.uber.org/zap"
)

const (
	loggerKey = utils.ContextKey("logger")
)

func NewLogger(debug bool) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func AddLoggerToContext(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger)
	if !ok {
		return zap.NewNop().Sugar()
	}
	return logger
}
