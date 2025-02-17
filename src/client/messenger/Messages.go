package messenger

import "crypto/rsa"

type HandshakeMessage struct {
	Infohash  [20]byte
	Id        string
	PublicKey *rsa.PublicKey
}

type ChokeMessage struct{}

type UnchokeMessage struct{}

type InterestedMessage struct{}

type NotInterestedMessage struct{}

type HaveMessage struct {
	Index int
}

type BitfieldMessage struct {
	Bitfield []bool
}

type RequestMessage struct {
	Index  int
	Offset int
	Length int
}

type PieceMessage struct {
	Index  int
	Offset int
	Bytes  []byte
}

type CancelMessage struct {
	RequestMessage
}
