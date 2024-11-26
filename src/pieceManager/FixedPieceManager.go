package pieceManager

import (
	"bittorrent/common"
	"errors"
)

type fixedPieceManager struct {
	chunks          map[int][]bool
	uncheckedChunks map[int]int

	chunkLength          int
	lastChunkLength      int
	lastPieceChunkLength int
}

func New(length int, pieceLength int, chunkLength int) PieceManager {
	manager := fixedPieceManager{}
	manager.chunkLength = chunkLength
	manager.chunks = map[int][]bool{}
	manager.uncheckedChunks = map[int]int{}

	// Assuming the pieceLength >= chunkLength
	chunksPerPiece := pieceLength / chunkLength
	manager.lastChunkLength = pieceLength % chunkLength
	if manager.lastChunkLength != 0 {
		chunksPerPiece += 1
	}

	totalPieces := common.GetTotalPieces(length, pieceLength)
	for i := range totalPieces {
		// Handle last piece bytes since they might be truncated
		if i == totalPieces-1 && length%pieceLength != 0 {
			lastPieceBytes := length % pieceLength

			var lastPieceChunks int

			if chunkLength > lastPieceBytes {
				lastPieceChunks = 1
			} else {
				lastPieceChunks = lastPieceBytes / chunkLength
				manager.lastPieceChunkLength = lastPieceBytes % chunkLength
				if manager.lastPieceChunkLength != 0 {
					lastPieceChunks += 1
				}
			}

			chunksPerPiece = lastPieceChunks
		}

		manager.chunks[i] = make([]bool, chunksPerPiece)
		for j := range chunksPerPiece {
			manager.chunks[i][j] = true
		}
		manager.uncheckedChunks[i] = chunksPerPiece
	}

	return &manager
}

func (manager *fixedPieceManager) ChunkLength() int {
	return manager.chunkLength
}

func (manager *fixedPieceManager) CheckChunk(index int, offset int) bool {
	chunkIndex := offset / manager.chunkLength
	if manager.chunks[index][chunkIndex] {
		return manager.uncheckedChunks[index] == 0
	}

	manager.chunks[index][chunkIndex] = true
	manager.uncheckedChunks[index] -= 1
	return manager.uncheckedChunks[index] == 0
}

func (manager *fixedPieceManager) UncheckPiece(index int) {
	chunks := len(manager.chunks[index])
	manager.chunks[index] = make([]bool, chunks)
	manager.uncheckedChunks[index] = chunks
}

func (manager *fixedPieceManager) GetUncheckedPieces() []int {
	uncheckedPieces := []int{}
	for index, uncheckedChunks := range manager.uncheckedChunks {
		if uncheckedChunks > 0 {
			uncheckedPieces = append(uncheckedPieces, index)
		}
	}
	return uncheckedPieces
}

func (manager *fixedPieceManager) GetUncheckedChunk(index int) (int, int, int, error) {
	for chunkIndex, isCheckedChunk := range manager.chunks[index] {
		if !isCheckedChunk {
			isLastPiece := index == len(manager.chunks)-1
			isLastChunk := chunkIndex == len(manager.chunks[index])-1

			offset := chunkIndex * manager.chunkLength
			length := manager.chunkLength

			if isLastChunk && manager.lastChunkLength > 0 {
				length = manager.lastChunkLength
			}

			if isLastPiece && isLastChunk {
				length = manager.lastPieceChunkLength
			}

			return index, offset, length, nil
		}
	}

	return -1, -1, -1, errors.New("all chunks are already checked")
}
