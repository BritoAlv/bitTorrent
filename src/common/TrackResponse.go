package common

type TrackResponse struct {
	FailureReason string             `bencode:"failure reason"`
	Interval      int                `bencode:"interval"`
	Peers         map[string]Address `bencode:"peers"`
}
