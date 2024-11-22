package common

import "os"

type Logger struct {
	FileName string
}

func NewLogger(fileID string) *Logger {
	_, _ = os.OpenFile(fileID, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	return &Logger{FileName: fileID}
}

func (l *Logger) WriteToFile(content string) {
	file, err := os.OpenFile(l.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(content + "\n")
	if err != nil {
		panic(err)
	}
}
