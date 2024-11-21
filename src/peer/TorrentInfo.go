package peer

type TorrentInfo struct {
	Hash        []byte // Info-hash of the .torrent file
	Name        string // Suggested name to save the file (or directory) as. It is purely advisory.
	PieceLength uint32 // The number of bytes in each piece the file is split into
	Pieces      []byte
	Length      uint32     // It's present when the download represents a single file
	Files       []FileInfo // It's present when the download represents a directory (if not Files == nil). It represents a set of files which go in a directory structure
}
