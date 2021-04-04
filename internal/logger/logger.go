package logger

import (
	"github.com/op/go-logging"
)

func New(module string) *logging.Logger {
	log := logging.MustGetLogger(module)
	return log
}
