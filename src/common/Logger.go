package common

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var LogsPath = "./logs/"

type Logger struct {
	DefaultPath string
	FileName    string
	Prefix      string
	lock        sync.Mutex
}

func SetLogDirectoryPath(name string) {
	curTime := time.Now()
	formattedTime := curTime.Format("15:04:05")
	LogsPath = "./logs/" + name + formattedTime + "/"
}

func NewLogger(fileID string) *Logger {
	defaultPath := LogsPath
	err := os.MkdirAll(defaultPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	_, _ = os.OpenFile(defaultPath+fileID, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	return &Logger{
		DefaultPath: defaultPath,
		FileName:    fileID,
		Prefix:      "",
		lock:        sync.Mutex{},
	}
}

func (l *Logger) SetPrefix(prefix string) {
	l.Prefix = prefix
}

func (l *Logger) WriteToFileError(format string, args ...interface{}) {
	format = "[ERROR] " + format
	l.WriteToFileOK(format, args)
}

func (l *Logger) WriteToFileOK(format string, args ...interface{}) {
	content := fmt.Sprintf(format, args...)
	l.lock.Lock()
	defer l.lock.Unlock()
	file, err := os.OpenFile(l.DefaultPath+l.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(content + "\n")
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		return
	}
}
