package pieceManager

type PieceManager interface {
	ChunkLength() int
	CheckChunk(index int, offset int) bool
	UncheckPiece(index int)
	GetUncheckedPieces() []int
	GetUncheckedChunk(index int) (int, int, int, error) // Returns (index, offset, length)
}
