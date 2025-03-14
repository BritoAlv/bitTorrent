package pieceManager

import (
	"bittorrent/common"
	"sync"
)

type fixedPieceManager struct {
	chunks          map[int][]bool
	uncheckedChunks map[int]int

	chunkLength          int
	lastChunkLength      int
	lastPieceChunkLength int

	mutex *sync.Mutex
}

func New(length int, pieceLength int, chunkLength int) PieceManager {
	manager := fixedPieceManager{}
	manager.chunkLength = chunkLength
	manager.chunks = map[int][]bool{}
	manager.uncheckedChunks = map[int]int{}
	manager.mutex = new(sync.Mutex)

	// Assuming the pieceLength >= chunkLength
	chunksPerPiece := pieceLength / chunkLength
	manager.lastChunkLength = pieceLength % chunkLength
	if manager.lastChunkLength != 0 {
		chunksPerPiece += 1
	}

	totalPieces := common.GetTotalPieces(length, pieceLength)
	for i := range totalPieces {
		// HandleNotification last piece bytes since they might be truncated
		if i == totalPieces-1 && length%pieceLength != 0 {
			lastPieceBytes := length % pieceLength

			var lastPieceChunks int

			if chunkLength > lastPieceBytes {
				lastPieceChunks = 1
				manager.lastPieceChunkLength = lastPieceBytes
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
		manager.uncheckedChunks[i] = 0
	}

	return &manager
}

func (manager *fixedPieceManager) ChunkLength() int {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	return manager.chunkLength
}

func (manager *fixedPieceManager) Bitfield() []bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	bitfield := make([]bool, len(manager.uncheckedChunks))

	for index, uncheckedChunks := range manager.uncheckedChunks {
		if uncheckedChunks == 0 {
			bitfield[index] = true
		} else {
			bitfield[index] = false
		}
	}

	return bitfield
}

func (manager *fixedPieceManager) VerifyPiece(index int) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	return manager.uncheckedChunks[index] == 0
}

func (manager *fixedPieceManager) VerifyChunk(index int, offset int) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	chunkIndex := offset / manager.chunkLength
	return manager.chunks[index][chunkIndex]
}

func (manager *fixedPieceManager) CheckChunk(index int, offset int) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	chunkIndex := offset / manager.chunkLength
	if manager.chunks[index][chunkIndex] {
		return manager.uncheckedChunks[index] == 0
	}

	manager.chunks[index][chunkIndex] = true
	manager.uncheckedChunks[index] -= 1
	return manager.uncheckedChunks[index] == 0
}

func (manager *fixedPieceManager) UncheckPiece(index int) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	chunks := len(manager.chunks[index])
	manager.chunks[index] = make([]bool, chunks)
	manager.uncheckedChunks[index] = chunks
}

func (manager *fixedPieceManager) GetUncheckedPieces() []int {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	uncheckedPieces := []int{}
	for index, uncheckedChunks := range manager.uncheckedChunks {
		if uncheckedChunks > 0 {
			uncheckedPieces = append(uncheckedPieces, index)
		}
	}
	return uncheckedPieces
}

func (manager *fixedPieceManager) GetUncheckedChunks(index int, n int) [][3]int {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	uncheckedChunks := [][3]int{}

	count := 0
	for chunkIndex, isCheckedChunk := range manager.chunks[index] {
		if !isCheckedChunk {
			isLastPiece := index == len(manager.chunks)-1
			isLastChunk := chunkIndex == len(manager.chunks[index])-1

			offset := chunkIndex * manager.chunkLength
			length := manager.chunkLength

			if isLastChunk && manager.lastChunkLength > 0 {
				length = manager.lastChunkLength
			}

			if isLastPiece && isLastChunk && manager.lastPieceChunkLength > 0 {
				length = manager.lastPieceChunkLength
			}

			uncheckedChunks = append(uncheckedChunks, [3]int{index, offset, length})
			count++
		}

		if count > n {
			break
		}
	}

	return uncheckedChunks
}
