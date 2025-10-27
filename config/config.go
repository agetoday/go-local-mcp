package config

import (
	"log"
	"os"
	"time"
)

// 日志配置
func init() {
	rLog()
	log.Println("Starting config initialization...")

}

func rLog() {
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "logs"
	}
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, 0755); err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}
	}
	dateStamp := time.Now().Format("2006-01-02")
	logFileName := logDir + "/" + dateStamp + ".log"
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
