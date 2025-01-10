package utils

import "strings"

func OkOrPanic(err error, msg ...string) {
	if err != nil {
		errInfo := strings.Join(msg, " ") + err.Error()
		panic(errInfo)
	}
}
