package common

import "strconv"

type OfficialTrackAddress struct {
	Ip   string `bencode:"ip"`
	Port int    `bencode:"port"`
}

type OfficialTrackResponse struct {
	FailureReason string                 `bencode:"failure reason,omitempty"`
	Interval      int                    `bencode:"interval"`
	Peers         []OfficialTrackAddress `bencode:"peers"`
}

func BuildOfficialResponse(response TrackResponse) OfficialTrackResponse {
	var officialResponse OfficialTrackResponse
	officialResponse.FailureReason = response.FailureReason
	officialResponse.Interval = response.Interval
	officialResponse.Peers = make([]OfficialTrackAddress, 0)
	for _, address := range response.Peers {
		port, _ := strconv.Atoi(address.Port)
		officialResponse.Peers = append(officialResponse.Peers, OfficialTrackAddress{
			Ip:   address.Ip,
			Port: port,
		})
	}
	return officialResponse
}
