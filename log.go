package log

import "github.com/weeon/contract"

func GetLogger() contract.Logger {
	return _logger
}
