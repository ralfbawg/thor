package logging

import (
	"fmt"
	"time"
	"os"
	"github.com/panjf2000/ants"
)

const (
	errorLevel = "error"
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
)

var defaultLogFile = os.TempDir() + string(os.PathSeparator) + "thor.log"

type LogManager struct {
	logFileObj  *os.File
	logFilePath string
	log         chan string
	logLevel    string
}

var LogManagerInst = &LogManager{logLevel: debugLevel, log: make(chan string, 100),}

func Init(level string, logPath string) {
	if LogManagerInst != nil {
		if level != "" {
			LogManagerInst.logLevel = level
		}
		if logPath != "" {
			LogManagerInst.logFilePath = logPath
		}
	}
	fileObj, err := os.Create(logPath)
	if err == nil {
		LogManagerInst.logFileObj = fileObj
	} else {
		fmt.Println("日志文件初始化失败")
	}
	ants.Submit(LogManagerInst.run)
}

func (l *LogManager) SetLogPath(path string) {
	fileObj, err := os.Create(path)
	if err == nil {
		LogManagerInst.logFileObj = fileObj
	} else {
		fmt.Println("日志文件切换失败")
	}
}
func (l *LogManager) run() {
	for {
		select {
		case log := <-l.log:
			fmt.Println(log)
			l.logFileObj.WriteString(log + "\n")
		}
	}
}
func (l *LogManager) SetLogLevel(level string) {
	l.logLevel = level
}

func Debug(format string, values ...interface{}) {
	switch LogManagerInst.logLevel {
	case debugLevel:
		printLog("DEBUG", format, values...)
	}
}

func Info(format string, values ...interface{}) {
	switch LogManagerInst.logLevel {
	case debugLevel, errorLevel, warnLevel, infoLevel:
		printLog("INFO", format, values...)
	}

}

func Error(format string, values ...interface{}) {
	printLog("ERROR", format, values...)

}

func Warning(format string, values ...interface{}) {
	printLog("WARINING", format, values...)
}

func printLog(keyword string, format string, values ...interface{}) {
	if LogManagerInst == nil {
		Init(debugLevel, defaultLogFile)
	}
	time_string := time.Now().Format("2006-01-02 15:04:05")
	log_string := "[" + time_string + "] " + keyword + " " + fmt.Sprintf(format, values...)
	LogManagerInst.log <- log_string
}
