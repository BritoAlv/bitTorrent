package common

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/zeebo/bencode"
)

type Torrent struct {
	Announce    string // Url of the tracker
	InfoHash    []byte // Info-hash of the .torrent file
	Name        string // Suggested name to save the file (or directory) as. It is purely advisory.
	PieceLength int64  // The number of bytes in each piece the file is split into
	Pieces      []byte
	Length      int64      // It's present when the download represents a single file
	Files       []FileInfo // It's present when the download represents a directory (if not Files == nil). It represents a set of files which go in a directory structure
}

// Torrent keys
const _ANNOUNCE = "announce"

// const _CREATED_BY = "created by"
// const _CREATION_DATE = "creation date"
const _INFO = "info"
const _NAME = "name"
const _LENGTH = "length"
const _PIECE_LENGTH = "piece length"
const _PIECES = "pieces"

// const _ANNOUNCE_LIST = "announce-list"
// const _COMMENT = "comment"

// File access permission
const _PERMISSION = 0644 // User can read/write. Groups can only read

func ParseTorrentFile(fileName string) (Torrent, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR, _PERMISSION)

	if err != nil {
		return Torrent{}, fmt.Errorf("an error occurred while parsing the torrent file: %w", err)
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return Torrent{}, fmt.Errorf("an error occurred while parsing the torrent file: %w", err)
	}

	fileSize := fileInfo.Size()

	buffer := make([]byte, fileSize)
	bytesRead, err := file.Read(buffer)

	if err != nil {
		return Torrent{}, fmt.Errorf("an error occurred while parsing the torrent file: %w", err)
	}

	bytes := buffer[:bytesRead]

	var torrentDic map[string]interface{}
	err = bencode.DecodeBytes(bytes, &torrentDic)

	if err != nil {
		return Torrent{}, fmt.Errorf("an error occurred while parsing the torrent file: %w", err)
	}

	torrent, err := extractTorrent(torrentDic)

	if err != nil {
		return Torrent{}, fmt.Errorf("an error occurred while parsing the torrent file: %w", err)
	}

	torrentInfo, err := CastTo[map[string]interface{}](torrentDic[_INFO])

	if err != nil {
		return Torrent{}, fmt.Errorf("an error occurred while parsing the torrent file: %w", err)
	}

	torrentInfoBytes, err := bencode.EncodeBytes(torrentInfo)

	if err != nil {
		return Torrent{}, fmt.Errorf("an error occurred while parsing the torrent file: %w", err)
	}

	infoHash := sha1.Sum(torrentInfoBytes)

	torrent.InfoHash = infoHash[:]

	return torrent, nil
}

func extractTorrent(torrentDic map[string]interface{}) (Torrent, error) {
	announce, err := CastTo[string](torrentDic[_ANNOUNCE])

	if err != nil {
		return Torrent{}, err
	}

	torrentInfo, err := CastTo[map[string]interface{}](torrentDic[_INFO])

	if err != nil {
		return Torrent{}, err
	}

	name, err := CastTo[string](torrentInfo[_NAME])

	if err != nil {
		return Torrent{}, err
	}

	length, err := CastTo[int64](torrentInfo[_LENGTH])

	if err != nil {
		return Torrent{}, err
	}

	pieceLength, err := CastTo[int64](torrentInfo[_PIECE_LENGTH])

	if err != nil {
		return Torrent{}, err
	}

	pieces, err := CastTo[string](torrentInfo[_PIECES])

	if err != nil {
		return Torrent{}, err
	}

	return Torrent{
		Announce:    announce,
		InfoHash:    nil,
		Name:        name,
		PieceLength: pieceLength,
		Pieces:      []byte(pieces),
		Length:      length,
		Files:       nil,
	}, nil
}
