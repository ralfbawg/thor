package logging

import (
	"fmt"
	"time"
)

const (
	errorLevel = "error"
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
)

var currentLogLevel = infoLevel

func Init(level string) {
	currentLogLevel = level
}

func SetLogLevel(level string) {
	currentLogLevel = level
}

func Debug(format string, values ...interface{}) {
	switch currentLogLevel {
	case debugLevel:
		time_string := time.Now().Format("2006-01-02 15:04:05")
		log_string := "[" + time_string + "] " + "DEBUG" + " " + fmt.Sprintf(format, values...)
		fmt.Println(log_string)
	}
}

func Info(format string, values ...interface{}) {
	switch currentLogLevel {
	case errorLevel, warnLevel, infoLevel:
		time_string := time.Now().Format("2006-01-02 15:04:05")
		log_string := "[" + time_string + "] " + "INFO" + " " + fmt.Sprintf(format, values...)
		fmt.Println(log_string)
	}

}

func Error(format string, values ...interface{}) {
	time_string := time.Now().Format("2006-01-02 15:04:05")
	log_string := "[" + time_string + "] " + "ERROR" + " " + fmt.Sprintf(format, values...)
	fmt.Println(log_string)
}

func Warning(format string, values ...interface{}) {
	time_string := time.Now().Format("2006-01-02 15:04:05")
	log_string := "[" + time_string + "] " + "WARNING" + " " + fmt.Sprintf(format, values...)
	fmt.Println(log_string)
}
