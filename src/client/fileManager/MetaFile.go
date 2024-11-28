package fileManager

import "os"

type metaFile struct {
	Index         int
	FileReference *os.File
	Length        int
}
