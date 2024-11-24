package fileManager

import (
	"bittorrent/common"
	"errors"
	"fmt"
	"os"
	"path"
)

type ConcurrentFileManager struct {
	files []metaFile
}

// **Public methods

func NewConcurrentFileManager(fileInfos []common.FileInfo) (*ConcurrentFileManager, error) {
	fileManager := ConcurrentFileManager{
		files: []metaFile{},
	}
	fileIndex := 0

	for _, info := range fileInfos {
		file, err := extractFile(info.Path)
		if err != nil {
			return nil, err
		}

		fileManager.files = append(fileManager.files, metaFile{
			Index:         fileIndex,
			FileReference: file,
			Length:        info.Length,
		})
		fileIndex += int(info.Length)
	}

	return &fileManager, nil
}

func (fileWriter *ConcurrentFileManager) Write(start int, bytes *[]byte) (bool, error) {
	return true, errors.New("not implemented")
}

func (fileWriter *ConcurrentFileManager) Read(start int, length int) ([]byte, error) {
	bytes := []byte{}

	end := start + length - 1
	for _, metaFile := range fileWriter.files {
		fileStart := metaFile.Index
		fileEnd := fileStart + metaFile.Length - 1
		var properStart, properEnd int

		if start > fileEnd {
			continue
		}

		if end < fileStart {
			break
		}

		if start >= fileStart {
			properStart = start - fileStart
		} else {
			properStart = 0
		}

		if end >= fileEnd {
			properEnd = metaFile.Length - 1
		} else {
			properEnd = end - fileStart
		}

		properLength := properEnd - properStart + 1
		targetBytes := make([]byte, properLength)
		_, err := metaFile.FileReference.ReadAt(targetBytes, int64(properStart))
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, targetBytes...)
	}

	return bytes, nil
}

//** Private methods

func extractFile(fileName string) (*os.File, error) {
	directory, _ := path.Split(fileName)

	err := os.MkdirAll(directory, common.ALL_RW_PERMISSION)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, common.ALL_RW_PERMISSION)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return file, nil
}
