package tracker

import "bittorrent/common"

type CentralizedTracker struct {
	Address common.Address
}

func (tracker CentralizedTracker) Track(request TrackRequest) (TrackResponse, error) {
	// Mocking a response
	response := TrackResponse{
		FailureReason: "",
		Interval:      10,
		Peers: map[common.Address]string{
			{Ip: "localhost", Port: "8085"}: "nature",
			{Ip: "localhost", Port: "8090"}: "ocean",
		},
	}
	return response, nil
}
