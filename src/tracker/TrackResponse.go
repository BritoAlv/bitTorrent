package tracker

import "bittorrent/common"

type TrackResponse struct {
	FailureReason string
	Interval      int
	Peers         map[string]common.Address // Peers is an <Id, Address> dictionary
}
