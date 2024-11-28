package peer

import (
	"bittorrent/common"
	"net"
)

type trackNotification struct {
	Response   common.TrackResponse
	Successful bool
}

type downloadNotification struct{}

type killNotification struct{}

type addPeerNotification struct {
	PeerId     string
	Connection net.Conn
}

type removePeerNotification struct {
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
