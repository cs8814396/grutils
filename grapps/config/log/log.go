package log

import (
	"github.com/gdgrc/grutils/grapps/config"
)

func Fine(format string, v ...interface{}) error {

	return config.DefaultLogger.OutPut(1, 2, "[Fine] ", format, v...)
}

func Debug(format string, v ...interface{}) error {

	return config.DefaultLogger.OutPut(2, 2, "[Debug] ", format, v...)
}

func Info(format string, v ...interface{}) error {
	return config.DefaultLogger.OutPut(3, 2, "[Info] ", format, v...)
}
func Warn(format string, v ...interface{}) error {
	return config.DefaultLogger.OutPut(4, 2, "[Warn] ", format, v...)
}

func Error(format string, v ...interface{}) (err error) {

	return config.DefaultLogger.OutPut(5, 2, "[Error] ", format, v...)

}
