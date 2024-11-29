package tracker

import (
	"bittorrent/common"
	"fmt"
	"io"
	"net/http"
)

type CentralizedHttpTracker struct {
	Url string
}


func (tracker CentralizedHttpTracker) Track(request common.TrackRequest) (common.TrackResponse, error) {
	UrlSend, err := common.BuildHttpUrl(tracker.Url, request)

	if err != nil {
		return common.TrackResponse{}, fmt.Errorf("an error occurred while building the url to contact the tracker : %w", err)
	}
	// Send GET request
	httpResponse, err := http.Get(UrlSend)
	if err != nil {
		return common.TrackResponse{}, fmt.Errorf("an error occurred while contacting the tracker: %w", err)
	}

	bytes, err := io.ReadAll(httpResponse.Body)
	if err != nil{
		return common.TrackResponse{}, fmt.Errorf("an error ocurred while reading the body of the response from the tracker: %w", err)
	}

	response, err := common.DecodeTrackerResponse(bytes)
	if err != nil {
		return common.TrackResponse{}, fmt.Errorf("an error occurred while decoding the response from the tracker: %w", err)
	}
	return response, nil
}