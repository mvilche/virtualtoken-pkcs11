package main

import (
	"log"
	"os"
	"strings"
	"time"
)

var InfoLog *log.Logger
var ExecutionLog *log.
var ErrorLog *log.Logger

func createLogFile() error {
	logdir := "logs/"
	date := strings.TrimSpace(time.Now().Format("2006-01-02"))
	logfile := logdir + date + ".log"

	if _, err := os.Stat(logdir); os.IsNotExist(err) {
		err := os.Mkdir(logdir, os.ModePerm)
		if err != nil {

			return err
		}
	}
	openLogfile, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	InfoLog = log.New(openLogfile, "INFO:\t", log.Ldate|log.Ltime)
	ExecutionLog = log.New(openLogfile, "RUN:\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(openLogfile, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
