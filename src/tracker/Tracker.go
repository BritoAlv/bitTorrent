package tracker

import "bittorrent/common"

type Tracker interface {
	Track(request common.TrackRequest) (response common.TrackResponse, err error)
}
