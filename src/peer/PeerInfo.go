package peer

import "net"

type PeerInfo struct {
	Connection net.Conn
	Bitfield   []bool
	IsChoker   bool
	IsChoked   bool
}
