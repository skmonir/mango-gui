package logger

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func Info(message string) {
	fileInfo, err := getLogFile()
	if err != nil {
		log.Fatal(err)
		return
	}
	infoLog := log.New(fileInfo, "[info] ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	infoLog.Println(message)
}

func Error(message string) {
	fileInfo, err := getLogFile()
	if err != nil {
		log.Fatal(err)
		return
	}
	infoLog := log.New(fileInfo, "[error] ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	infoLog.Println(message)
}

func getLogFile() (*os.File, error) {
	currentTime := time.Now()
	logDir := filepath.Join(utils.GetAppHomeDirectoryPath(), "logs")
	filename := currentTime.Format("2006-01-02") + ".log"
	logfile := filepath.Join(logDir, filename)

	if err := utils.CreateFile(logDir, filename); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

func InvalidateLogFiles() {
	logDir := filepath.Join(utils.GetAppHomeDirectoryPath(), "logs")
	logFiles := utils.GetFileNamesInDirectory(logDir)
	currentTime, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

	for _, filename := range logFiles {
		basename := strings.TrimRight(filename, filepath.Ext(filename))
		date, _ := time.Parse("2006-01-02", basename)
		nxtDate := date.AddDate(0, 0, 9)
		if nxtDate.Before(currentTime) {
			if err := utils.RemoveFile(filepath.Join(logDir, filename)); err != nil {
				Error("Error while removing " + filename)
			}
		}
	}
}
