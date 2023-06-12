package logtool

import (
	"go.uber.org/zap"
)

type ZapPrintfLogger struct {
	logger *zap.SugaredLogger
}

func (zl *ZapPrintfLogger) Printf(format string, args ...interface{}) {
	zl.logger.Infof(format, args...)
}
