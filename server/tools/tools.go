package tools

import (
	"../logger"
)

func FailOnWarning(err error, msg string) {
	if err != nil {
		logger.WarningPrintf("%s: %s", msg, err)
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		logger.ErrorPrintf("%s: %s", msg, err)
	}
}
