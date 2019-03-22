package app

import (
	"github.com/astaxie/beego/validation"
	"go-gin-demo/pkg/logging"
	"go.uber.org/zap"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.HTTPLogger.Error("validation error", zap.Error(err))
	}

	return
}
