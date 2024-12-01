package common

type TrackResponse struct {
	FailureReason string `bencode:"failure reason"`	
	Interval      int `bencode:"interal"`
	Peers         map[string]Address `bencode:"peers"`
}