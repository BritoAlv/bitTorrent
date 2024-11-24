package fileManager

type FileManager interface {
	Write(start int, bytes *[]byte) (bool, error)
	Read(start int, length int) ([]byte, error)
}
