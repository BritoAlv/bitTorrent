package common

import (
	"crypto/sha1"
	"fmt"
	"math"
	"os"

	"github.com/zeebo/bencode"
)

type Torrent struct {
	Announce    string   // Location of the tracker
	InfoHash    [20]byte // Info-hash of the .torrent file
	Name        string   // Suggested name to save the file (or directory) as. It is purely advisory.
	PieceLength int64    // The number of bytes in each piece the file is split into
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

// TODO: Check about this file access permission and what they actually mean
// File access permission
const USER_RW_PERMISSION = 0644 // User can read/write. Groups can only read
const ALL_RW_PERMISSION = 0777  // All can read/write

func ParseTorrentFile(fileName string) (Torrent, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR, USER_RW_PERMISSION)

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

	torrent.InfoHash = infoHash

	return torrent, nil
}

func CreateTorrentFile(path string, announceUrl string, isDirectory bool) error {
	file, err := os.OpenFile(path, os.O_RDWR, ALL_RW_PERMISSION)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	name := file.Name()
	pieceLength := int(math.Pow(2, 18))
	length := fileInfo.Size()
	totalPieces := GetTotalPieces(int(length), pieceLength)

	pieces := make([]byte, 0, totalPieces*20)
	buffer := make([]byte, pieceLength)

	start := 0
	for start < int(length) {
		bufferLength := len(buffer)
		end := start + bufferLength - 1

		// Check size of the last buffer
		if end >= int(length) {
			buffer = make([]byte, length-int64(start))
			bufferLength = len(buffer)
		}

		_, err := file.ReadAt(buffer, int64(start))
		if err != nil {
			return err
		}

		hashedBuffer := sha1.Sum(buffer)
		pieces = append(pieces, hashedBuffer[0:]...)

		start += bufferLength
	}

	torrentInf := map[string]interface{}{
		_NAME:         name,
		_PIECE_LENGTH: pieceLength,
		_PIECES:       pieces,
		_LENGTH:       length,
	}

	torrent := map[string]interface{}{
		_ANNOUNCE: announceUrl,
		_INFO:     torrentInf,
	}

	encodedTorrent, err := bencode.EncodeBytes(torrent)
	if err != nil {
		return err
	}

	torrentPath := "./" + name + ".torrent"
	torrentFile, err := os.OpenFile(torrentPath, os.O_CREATE|os.O_RDWR, ALL_RW_PERMISSION)
	if err != nil {
		return err
	}

	err = ReliableWrite(torrentFile, encodedTorrent)
	if err != nil {
		return err
	}

	return nil
}

func GetTotalPieces(length int, pieceLength int) int {
	var totalPieces int

	if length <= pieceLength {
		totalPieces = 1
	} else {
		totalPieces = length / pieceLength
		if length%pieceLength != 0 {
			totalPieces += 1
		}
	}
	return totalPieces
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
		InfoHash:    [20]byte{},
		Name:        name,
		PieceLength: pieceLength,
		Pieces:      []byte(pieces),
		Length:      length,
		Files:       nil,
	}, nil
}
