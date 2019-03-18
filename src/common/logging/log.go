package logging

import (
	"fmt"
	"time"
	"config"
)

const (
	errorLevel = "error"
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
)

func init() {
}

func Debug(format string, values ...interface{}) {
	switch config.ConfigStore.Log.Level {
	case debugLevel:
		time_string := time.Now().Format("2006-01-02 15:04:05")
		log_string := "[" + time_string + "] " + "INFO" + " " + fmt.Sprintf(format, values...)
		fmt.Println(log_string)
	}
}

func Info(format string, values ...interface{}) {
	switch config.ConfigStore.Log.Level {
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
