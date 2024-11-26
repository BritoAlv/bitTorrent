package fileManager

import (
	"bittorrent/common"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
)

type ConcurrentFileManager struct {
	files           []metaFile
	mutex           sync.Mutex
	lockedIntervals map[[2]int]bool
}

// **Public methods

func New(fileInfos []common.FileInfo) (FileManager, error) {
	fileManager := ConcurrentFileManager{
		files:           []metaFile{},
		mutex:           sync.Mutex{},
		lockedIntervals: map[[2]int]bool{},
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

func (fileWriter *ConcurrentFileManager) Write(start int, bytes *[]byte) error {
	end := start + len(*bytes) - 1
	targetInterval := [2]int{start, end}

	// Lock interval
	fileWriter.mutex.Lock()
	wasIntervalAdded := addInterval(fileWriter.lockedIntervals, targetInterval)
	fileWriter.mutex.Unlock()

	// Ensure interval is released after method terminates
	defer func() {
		fileWriter.mutex.Lock()
		delete(fileWriter.lockedIntervals, targetInterval)
		fileWriter.mutex.Unlock()
	}()

	if !wasIntervalAdded {
		return errors.New("interval is already taken")
	}

	bytesWritten := 0
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
		_, err := metaFile.FileReference.WriteAt((*bytes)[bytesWritten:bytesWritten+properLength], int64(properStart))
		if err != nil {
			return err
		}

		bytesWritten += properLength
	}

	return nil
}

func (fileWriter *ConcurrentFileManager) Read(start int, length int) ([]byte, error) {
	bytes := []byte{}
	end := start + length - 1
	targetInterval := [2]int{start, end}

	// Lock interval
	fileWriter.mutex.Lock()
	wasIntervalAdded := addInterval(fileWriter.lockedIntervals, targetInterval)
	fileWriter.mutex.Unlock()

	// Ensure interval is released after method terminates
	defer func() {
		fileWriter.mutex.Lock()
		delete(fileWriter.lockedIntervals, targetInterval)
		fileWriter.mutex.Unlock()
	}()

	if !wasIntervalAdded {
		return nil, errors.New("interval is already taken")
	}

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

func checkInterval(intervals map[[2]int]bool, targetInterval [2]int) bool {
	for interval := range intervals {
		if intervalContains(interval, targetInterval) {
			return false
		}
	}
	return true
}

func addInterval(intervals map[[2]int]bool, targetInterval [2]int) bool {
	if checkInterval(intervals, targetInterval) {
		intervals[targetInterval] = true
		return true
	}
	return false
}

// Checks if interval2 is contained in interval1
func intervalContains(interval1 [2]int, interval2 [2]int) bool {
	return interval1[0] <= interval2[0] && interval1[1] >= interval2[1]
}
