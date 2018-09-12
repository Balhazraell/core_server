package tools

import (
	"fmt"

	"../logger"
)

func FailOnWarning(err error, msg string) {
	if err != nil {
		logger.Warning.Printf("%s: %s", msg, err)
		fmt.Printf("%s: %s", msg, err)
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		logger.Error.Printf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
