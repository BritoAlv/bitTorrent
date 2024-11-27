package common

import "io"

func ReliableWrite(writer io.Writer, bytes []byte) error {
	totalWritten := 0

	for totalWritten < len(bytes) {
		bytesWritten, err := writer.Write(bytes[totalWritten:])

		if err != nil {
			return err
		}

		totalWritten += bytesWritten
	}

	return nil
}

func ReliableRead(reader io.Reader, bytesToRead int) ([]byte, error) {
	totalRead := 0
	bytes := []byte{}

	for totalRead < bytesToRead {
		bufferSize := bytesToRead - totalRead
		buffer := make([]byte, bufferSize)

		bytesRead, err := reader.Read(buffer)
		if err != nil {
			return nil, err
		}

		bytes = append(bytes, buffer[:bytesRead]...)
		totalRead += bytesRead
	}

	return bytes, nil
}
