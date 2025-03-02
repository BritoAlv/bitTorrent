package tracker

import "bittorrent/common"

type Tracker interface {
	Track(request common.TrackRequest) (response common.TrackResponse, err error)
}

func NewTracker(multicastUrl string) Tracker {
	return &multicastTracker{
		MulticastUrl: multicastUrl,
		ServerUrl:    "",
	}
}
