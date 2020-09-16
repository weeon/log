package log

import (
	"fmt"

	"github.com/weeon/contract"
	"go.uber.org/zap"
)

func GetLogger() contract.Logger {
	return _logger
}

func GetDefault() *Logger {
	return _logger
}

func SetDefault(l *Logger) {
	_logger = l
}

func FastInitFileLogger() {
	l, err := NewLogger("/app/log/normal.log", zap.DebugLevel)
	if err != nil {
		fmt.Println("Init file logger error ", err)
	} else {
		_logger = l
	}
}

func SetupStdoutLogger() {
	l, err := NewLogger("stdout", zap.DebugLevel)
	if err != nil {
		fmt.Println("Init stdout error ", err)
	} else {
		_logger = l
	}
}
