package log

import (
	"github.com/gdgrc/grutils/grapps/config"
)

func Fine(format string, v ...interface{}) error {

	return config.DefaultLogger.Fine(format, v...)
}

func Debug(format string, v ...interface{}) error {

	return config.DefaultLogger.Debug(format, v...)
}

func Info(format string, v ...interface{}) error {
	return config.DefaultLogger.Info(format, v...)
}
func Warn(format string, v ...interface{}) error {
	return config.DefaultLogger.Warn(format, v...)
}

func Error(format string, v ...interface{}) (err error) {

	return config.DefaultLogger.Error(format, v...)

}
