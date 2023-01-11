package elog

import (
	"github.com/jwping/logger"
)

var Log *logger.Logger

func init() {
	Log = logger.NewLogger(logger.Options{
		Lt:        logger.JSON,
		Level:     logger.Level(0),
		AddSource: true,
	})
}
