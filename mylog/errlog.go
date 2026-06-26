package mylog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// 独立的错误日志记录函数
func ErrLogToFile(level, message string) {
	if !enableFileLogGlobal {
		return
	}
	filename := getCurrentLogFilename()
	if err := os.MkdirAll(logPath, 0755); err != nil {
		fmt.Println("Error creating log directory:", err)
		return
	}
	filePath := filepath.Join(logPath, filename)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	logEntry := formatLogLine("ERROR", fmt.Sprintf("[%s] %s", level, message)) + "\n"
	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

// 独立的错误日志记录函数
func ErrInterfaceToFile(level, message interface{}) {
	if !enableFileLogGlobal {
		return
	}
	filename := getCurrentLogFilename()
	if err := os.MkdirAll(logPath, 0755); err != nil {
		fmt.Println("Error creating log directory:", err)
		return
	}
	filePath := filepath.Join(logPath, filename)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling data for log: %s", err)
		return
	}

	logEntry := formatLogLine("ERROR", fmt.Sprintf("[%s] %s", level, string(jsonData))) + "\n"
	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}
