package tracker

type Tracker interface {
	Track(request TrackRequest) (response TrackResponse, err error)
}
