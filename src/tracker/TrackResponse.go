package tracker

import "bittorrent/common"

type TrackResponse struct {
	FailureReason string
	Interval      int
	Peers         map[common.Address]string // Peers is an <Address, Id> dictionary
}
