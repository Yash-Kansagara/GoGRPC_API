package utils

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func ErrorHandler(err error, message string) error {
	if logger == nil {
		logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	logger.Println(message, err)
	return fmt.Errorf(message, err)
}
