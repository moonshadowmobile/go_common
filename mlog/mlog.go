package mlog

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var prefixMap map[*log.Logger]string
/* Creates a new Logger and returns a pointer to that Logger. */
func GetLogger(logName, prefix, logPath string) (*log.Logger, error) {
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, err
	}
	logFile, err := os.OpenFile(
		strings.Join([]string{logPath, logName}, ""),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0777,
	)

	if err != nil {
		return nil, err
	}

	logger := log.New(logFile, prefix, log.LstdFlags)
	if prefixMap == nil {
		prefixMap = make(map[*log.Logger]string)
	}
	prefixMap[logger] = prefix

	return logger, nil
}

func Println(logger *log.Logger, msg string) {
	fmt.Println(prefixMap[logger] + msg)
	logger.Println(msg)
}

func Fatalln(logger *log.Logger, msg string) {
	fmt.Println(prefixMap[logger] + msg)
	logger.Fatalln(msg)
}
