package testsuite

import (
	"github.com/tweethub/backend/pkg/log"
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	logger, err := log.New("test", "debug")
	if err != nil {
		panic(err)
	}
	return logger
}
