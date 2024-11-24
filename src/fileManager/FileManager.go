package fileManager

type FileManager interface {
	Write(start int, bytes *[]byte) error
	Read(start int, length int) ([]byte, error)
}
