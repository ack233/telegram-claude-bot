package logtool

import (
	"strings"

	"go.uber.org/zap"
)

func Fatalerror(_e error, s ...string) {
	if _e != nil {
		Logc.Fatal(strings.Join(s,""), zap.Error(_e))
	}
}

func Errorerror(_e error, s ...string) {
	if _e != nil {
		Logc.Error(strings.Join(s,""), zap.Error(_e))
	}
}
