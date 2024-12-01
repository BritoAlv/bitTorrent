package torrent

import (
	"bittorrent/common"
	"bittorrent/fileManager"
	"crypto/sha1"
	"fmt"
	"math"
	"os"
	"path"
	"strings"

	"github.com/zeebo/bencode"
)

type Torrent struct {
	Announce    string   // Location of the tracker
	InfoHash    [20]byte // Info-hash of the .torrent file
	Name        string   // Suggested name to save the file (or directory) as. It is purely advisory.
	PieceLength int64    // The number of bytes in each piece the file is split into
	Pieces      []byte
	Length      int64             // It's present when the download represents a single file
	Files       []common.FileInfo // It's present when the download represents a directory (if not Files == nil). It represents a set of files which go in a directory structure
}

// Torrent keys
const _ANNOUNCE = "announce"

// const _CREATED_BY = "created by"
// const _CREATION_DATE = "creation date"
const _INFO = "info"
const _NAME = "name"
const _LENGTH = "length"
const _FILES = "files"
const _PIECE_LENGTH = "piece length"
const _PIECES = "pieces"

// const _ANNOUNCE_LIST = "announce-list"
// const _COMMENT = "comment"

// TODO: Check about this file access permission and what they actually mean

func ParseTorrentFile(fileName string) (Torrent, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR, common.USER_RW_PERMISSION)

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

	torrentInfo, err := common.CastTo[map[string]interface{}](torrentDic[_INFO])

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

func CreateTorrentFile(targetPath string, torrentName string, announceUrl string) error {
	name := path.Base(targetPath)

	length, fileInfos, err := processPath(targetPath)
	if err != nil {
		return err
	}

	var isFile bool
	if fileInfos == nil {
		isFile = true
		fileInfos = []common.FileInfo{
			{
				Length: int(length),
				Path:   targetPath,
			},
		}
		targetPath = ""
	}

	manager, err := fileManager.New(targetPath, fileInfos)
	if err != nil {
		return err
	}

	pieceLength := int(math.Pow(2, 18))
	totalPieces := common.GetTotalPieces(int(length), pieceLength)

	pieces := make([]byte, 0, totalPieces*20)

	start := 0
	for start < int(length) {
		readLength := pieceLength
		end := start + readLength - 1

		// Check size of the last buffer
		if end >= int(length) {
			readLength = int(length) - start
		}

		bytes, err := manager.Read(start, readLength)
		if err != nil {
			return err
		}

		hashedBuffer := sha1.Sum(bytes)
		pieces = append(pieces, hashedBuffer[0:]...)

		start += readLength
	}

	var torrentInfo map[string]interface{}
	if isFile {
		torrentInfo = map[string]interface{}{
			_NAME:         name,
			_PIECE_LENGTH: pieceLength,
			_PIECES:       pieces,
			_LENGTH:       length,
		}
	} else {
		files := []map[string]interface{}{}

		for _, fileInfo := range fileInfos {
			files = append(files, map[string]interface{}{
				"length": fileInfo.Length,
				"path":   strings.Split(fileInfo.Path, "/")[1:],
			})
		}

		torrentInfo = map[string]interface{}{
			_NAME:         name,
			_PIECE_LENGTH: pieceLength,
			_PIECES:       pieces,
			_FILES:        files,
		}
	}

	torrent := map[string]interface{}{
		_ANNOUNCE: announceUrl,
		_INFO:     torrentInfo,
	}

	encodedTorrent, err := bencode.EncodeBytes(torrent)
	if err != nil {
		return err
	}

	torrentName = torrentName + ".torrent"

	// Remove file if there's already one named like that
	err = os.Remove(torrentName)
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		return err
	}

	torrentFile, err := os.OpenFile(torrentName, os.O_CREATE|os.O_RDWR, common.ALL_RW_PERMISSION)
	if err != nil {
		return err
	}

	err = common.ReliableWrite(torrentFile, encodedTorrent)
	if err != nil {
		return err
	}

	return nil
}

func extractTorrent(torrentDic map[string]interface{}) (Torrent, error) {
	announce, err := common.CastTo[string](torrentDic[_ANNOUNCE])

	if err != nil {
		return Torrent{}, err
	}

	torrentInfo, err := common.CastTo[map[string]interface{}](torrentDic[_INFO])

	if err != nil {
		return Torrent{}, err
	}

	name, err := common.CastTo[string](torrentInfo[_NAME])

	if err != nil {
		return Torrent{}, err
	}

	var length int64
	var files []common.FileInfo
	if _, contains := torrentInfo[_LENGTH]; contains {
		length, err = common.CastTo[int64](torrentInfo[_LENGTH])
		if err != nil {
			return Torrent{}, err
		}
	} else {
		files = []common.FileInfo{}
		rawFiles, err := common.CastTo[[]interface{}](torrentInfo[_FILES])
		if err != nil {
			return Torrent{}, nil
		}

		for _, rawFile := range rawFiles {
			fileInfo := common.FileInfo{}
			file, err := common.CastTo[map[string]interface{}](rawFile)
			if err != nil {
				return Torrent{}, nil
			}

			rawPathList, err := common.CastTo[[]interface{}](file["path"])
			if err != nil {
				return Torrent{}, nil
			}

			pathStr := ""
			for _, rawElem := range rawPathList {
				elem, err := common.CastTo[string](rawElem)
				if err != nil {
					return Torrent{}, nil
				}

				pathStr += "/" + elem
			}

			fileLength, err := common.CastTo[int64](file["length"])
			if err != nil {
				return Torrent{}, nil
			}

			fileInfo.Length = int(fileLength)
			fileInfo.Path = pathStr
			files = append(files, fileInfo)
		}
	}

	pieceLength, err := common.CastTo[int64](torrentInfo[_PIECE_LENGTH])

	if err != nil {
		return Torrent{}, err
	}

	pieces, err := common.CastTo[string](torrentInfo[_PIECES])

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
		Files:       files,
	}, nil
}

func processPath(targetPath string) (length int64, files []common.FileInfo, err error) {
	isDirectory := false
	file, err := os.OpenFile(targetPath, os.O_RDWR, common.ALL_RW_PERMISSION)
	if err != nil && strings.Contains(err.Error(), "is a directory") {
		isDirectory = true
	} else if err != nil {
		return 0, nil, err
	}

	if !isDirectory {
		fileInfo, err := file.Stat()
		if err != nil {
			return 0, nil, err
		}
		length = fileInfo.Size()
		return length, nil, nil
	} else {
		length := 0
		fileInfos, err := extractFileInfos(targetPath, true)
		if err != nil {
			return 0, nil, err
		}

		for _, fileInfo := range fileInfos {
			length += fileInfo.Length
		}
		return int64(length), fileInfos, nil
	}
}

func extractFileInfos(directory string, root bool) ([]common.FileInfo, error) {
	fileInfos := []common.FileInfo{}
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	baseName := "/" + path.Base(directory)

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			infos, err := extractFileInfos(directory+"/"+dirEntry.Name(), false)
			if err != nil {
				return nil, err
			}

			for _, info := range infos {
				if !root {
					info.Path = baseName + info.Path
				}
				fileInfos = append(fileInfos, info)
			}
		} else {
			info, err := dirEntry.Info()
			if err != nil {
				return nil, err
			}

			var pathStr string
			if !root {
				pathStr = baseName + "/" + info.Name()
			} else {
				pathStr = "/" + info.Name()
			}

			fileInfo := common.FileInfo{
				Length: int(info.Size()),
				Path:   pathStr,
			}
			fileInfos = append(fileInfos, fileInfo)
		}
	}

	return fileInfos, nil
}
