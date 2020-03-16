package grfile

import (
	"fmt"
	"github.com/gdgrc/grutils/grthird"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type XmlLogger struct {
	Level   int    `xml:"level"`
	LogFile string `xml:"logfile"`
}
type Logger struct {
	plogger         *log.Logger
	level           int
	filename        string
	logname         string
	alertOver       grthird.XmlAlertOver
	monitorSpace    float64
	lastMonitorTime time.Time
	monitorLock     sync.Mutex
}

func (mlog *Logger) SetMonitorAlertOver(ao grthird.XmlAlertOver) bool {
	mlog.alertOver = ao
	return true

}

func (mlog *Logger) CreateLogger(filename string, logname string, ilevel int) bool {

	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	//defer logFile.Close()
	if err != nil {
		log.Fatal("logger open file error ! filename: " + filename)
		return false
	}

	mlog.plogger = log.New(logFile, logname+"::", log.Ldate|log.Ltime)
	mlog.level = ilevel
	mlog.filename = filename
	mlog.logname = logname
	mlog.monitorSpace = 60.0
	mlog.lastMonitorTime = time.Now()

	return true

}

func (mlog *Logger) outPut(level int, preFix string, format string, v ...interface{}) error {

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("Logger Get FileLine err")
	}

	fileName := file

	a := strings.Split(file, "/")
	lenSplit := len(a)

	if lenSplit > 3 {
		fileName = a[lenSplit-2] + "/" + a[lenSplit-1]
	}

	lineStr := fmt.Sprintf("%d", line)

	if mlog.level <= level {

		mlog.plogger.Printf(fileName+":"+lineStr+" "+preFix+format, v...)
	}
	return nil
}
func (mlog *Logger) Fine(format string, v ...interface{}) error {

	return mlog.outPut(1, "[Fine] ", format, v...)
}

func (mlog *Logger) Debug(format string, v ...interface{}) error {

	return mlog.outPut(2, "[Debug] ", format, v...)
}

func (mlog *Logger) Info(format string, v ...interface{}) error {
	return mlog.outPut(3, "[Info] ", format, v...)
}
func (mlog *Logger) Warn(format string, v ...interface{}) error {
	return mlog.outPut(4, "[Warn] ", format, v...)
}

func (mlog *Logger) Error(format string, v ...interface{}) (err error) {
	err = mlog.outPut(5, "[Error] ", format, v...)
	if mlog.alertOver.Source != "" && mlog.alertOver.Receiver != "" {
		mlog.monitorLock.Lock()
		defer mlog.monitorLock.Unlock()
		nowTime := time.Now()
		deltaTime := nowTime.Sub(mlog.lastMonitorTime)
		if deltaTime.Seconds() >= mlog.monitorSpace {

			go grthird.AlertOverNotify(mlog.alertOver.Source, mlog.alertOver.Receiver, "ArmchairServer Error!", fmt.Sprintf(format, v...))
			mlog.lastMonitorTime = nowTime
		}

	}

	return err

}

func (mlog *Logger) GetLogger() (logger *log.Logger) {
	return mlog.plogger
}
func (mlog *Logger) SetLogger(logger *log.Logger) {
	mlog.plogger = logger

}

/*
func (mlog *Logger) CloneLogger() *Logger {

	var newLogger *Logger
	newLogger = new(Logger)
	newLogger.CreateLogger(mlog.filename, logname, ilevel)

	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	//defer logFile.Close()
	if err != nil {
		log.Fatal("logger open file error ! filename: " + filename)
		return false
	}

	mlog.plogger = log.New(logFile, logname+"::", log.Lshortfile|log.Ldate|log.Ltime)
	mlog.level = ilevel
	mlog.filename = filename
	mlog.logname = logname

}*/
