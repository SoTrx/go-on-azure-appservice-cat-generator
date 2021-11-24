package main

import (
	"log"
	"os"
)

type Logger struct {
	info *log.Logger
	warn *log.Logger
	err  *log.Logger
}

func MakeLogger() *Logger {
	logger := new(Logger)
	logger.info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.warn = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.err = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}
