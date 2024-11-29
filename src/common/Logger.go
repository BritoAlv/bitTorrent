package common

import (
	"fmt"
	"os"
)

type Logger struct {
	FileName string
}

func NewLogger(fileID string) *Logger {
	_, _ = os.OpenFile(fileID, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	return &Logger{FileName: fileID}
}

func (l *Logger) WriteToFileError(content string) {
	content = "Error (:" + content
	l.WriteToFileOK(content)
}

func (l *Logger) WriteToFileOK(content string) {
	file, err := os.OpenFile(l.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(content + "\n")
	fmt.Println(content)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		return
	}
}
