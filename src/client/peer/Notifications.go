package peer

import (
	"bittorrent/common"
	"net"
)

type killNotification struct{}

type trackNotification struct {
	Response   common.TrackResponse
	Successful bool
}

type downloadNotification struct{}

type writeNotification struct {
	Index  int
	Offset int
}

type pieceVerificationNotification struct {
	Index    int
	Verified bool
}

type addPeerNotification struct {
	PeerId     string
	Connection net.Conn
}

type removePeerNotification struct {
	PeerId string
}

type sendBitfieldNotification struct {
	PeerId string
}

type peerRequestNotification struct {
	PeerId string
	Index  int
	Offset int
	Length int
}

type peerCancelNotification struct {
	PeerId string
	Index  int
	Offset int
	Length int
}

type peerPieceNotification struct {
	PeerId string
	Index  int
	Offset int
	Bytes  []byte
}

type peerHaveNotification struct {
	PeerId string
	Index  int
}

type peerBitfieldNotification struct {
	PeerId   string
	Bitfield []bool
}

type peerChokeNotification struct {
	PeerId string
	Choke  bool
}

type peerInterestedNotification struct {
	PeerId     string
	Interested bool
}
