package peer

import (
	"crypto/rsa"
	"net"
)

type PeerInfo struct {
	Connection   net.Conn
	Bitfield     []bool
	IsInterested bool
	IsChoker     bool
	IsChoked     bool
	PublicKey    *rsa.PublicKey
}
