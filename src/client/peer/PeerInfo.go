package peer

import "net"

type PeerInfo struct {
	Connection   net.Conn
	Bitfield     []bool
	IsInterested bool
	IsChoker     bool
	IsChoked     bool
}
