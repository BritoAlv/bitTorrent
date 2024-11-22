package peer

import "net"

type PeerInfo struct {
	Id         string
	Connection net.Conn
	Bitfield   []bool
	IsChoker   bool
	IsChoked   bool
}
