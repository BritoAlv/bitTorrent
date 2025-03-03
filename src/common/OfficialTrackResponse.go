package common

type OfficialTrackAddress struct {
	Ip   string `bencode:"ip"`
	Port int    `bencode:"port"`
}

type OfficialTrackResponse struct {
	FailureReason string                 `bencode:"failure reason,omitempty"`
	Interval      int                    `bencode:"interval"`
	Peers         []OfficialTrackAddress `bencode:"peers"`
}
