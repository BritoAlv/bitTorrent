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
