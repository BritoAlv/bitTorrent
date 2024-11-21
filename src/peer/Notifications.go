package peer

import (
	"bittorrent/common"
	"bittorrent/tracker"
	"net"
)

type TrackerResponseNotification struct {
	Response tracker.TrackResponse
}

type DownloadNotification struct{}

type PeerDownNotification struct {
	Address common.Address
}

type PeerUpNotification struct {
	Address    common.Address
	Id         string
	Connection net.Conn
}

type KillNotification struct{}
