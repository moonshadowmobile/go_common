package go_common

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"strings"
	"time"
)

/* Creates a new Logger and returns a pointer to that Logger. */
func GetLogger(log_name string) (*log.Logger, error) {
	log_file, err := os.OpenFile(
		strings.Join([]string{C.LogPath, log_name}, ""),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0777,
	)

	if err != nil {
		return nil, err
	}

	logger := log.New(log_file, "", log.LstdFlags)

	return logger, nil
}

func Println(logger *log.Logger, msg string) {
	fmt.Println(time.Now().UTC().String() + ": " + msg)
	logger.Println(msg)
}

func Fatalln(logger *log.Logger, msg string) {
	fmt.Println(time.Now().UTC().String() + ": " + msg)
	logger.Fatalln(msg)
}

/* Set standard error log destination */
func setSysLog(logpath string) {
	err := os.MkdirAll(logpath, 0755)
	if err != nil {
		log.Fatalf("ERROR: Could not generate log paths. %s", err.Error())
	}

	error_log, err := os.OpenFile(
		strings.Join([]string{logpath, "sys_error.log"}, ""),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0777,
	)

	if err != nil {
		log.Fatalf("ERROR: Could not resolve path to logging destination. %s",
			err.Error())
	}
	log.SetOutput(io.MultiWriter(error_log, os.Stdout))
}
