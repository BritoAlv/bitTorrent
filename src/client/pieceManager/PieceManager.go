package pieceManager

type PieceManager interface {
	ChunkLength() int
	Bitfield() []bool

	VerifyChunk(index int, offset int) bool
	VerifyPiece(index int) bool

	CheckChunk(index int, offset int) bool
	UncheckPiece(index int)

	GetUncheckedPieces() []int
	GetUncheckedChunks(index int, n int) [][3]int // Gets at most n unchecked chunks
}
