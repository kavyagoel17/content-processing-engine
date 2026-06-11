package logger

import (
	"io"
	"log"
	"os"
)

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func Init() {
	logFile, err := os.OpenFile("appLogs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)

	InfoLog = log.New(multiWriter, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(multiWriter, "WARNING", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(multiWriter, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)
}
