package common

type TrackResponse struct {
	FailureReason string
	Interval      int
	Peers         map[string]Address // Peers is an <Id, Address> dictionary
}

type TrackFailureResponse struct {
	FailureReason string
}

type TrackSuccessResponse struct {
	Interval int
	Peers    map[string]Address
}