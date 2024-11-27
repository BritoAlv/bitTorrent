package tracker

import (
	"bittorrent/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type CentralizedTracker struct {
	Url string
}

const _INFO_HASH = "info_hash"
const _PEER_ID = "peer_id"
const _IP = "ip"
const _PORT = "port"
const _UPLOADED = "uploaded"
const _DOWNLOADED = "downloaded"
const _LEFT = "left"
const _EVENT = "event"

const _FAILURE_REASON = "failure_reason"
const _INTERVAL = "interval"
const _PEERS = "peers"

func (tracker CentralizedTracker) Track(request TrackRequest) (TrackResponse, error) {
	url, err := url.Parse(tracker.Url)

	if err != nil {
		return TrackResponse{}, fmt.Errorf("an error occurred while contacting tracker: %w", err)
	}

	values := url.Query()
	values.Set(_INFO_HASH, string(request.InfoHash))
	values.Set(_PEER_ID, request.PeerId)
	values.Set(_IP, request.Ip)
	values.Set(_PORT, request.Port)
	values.Set(_LEFT, strconv.Itoa(request.Left))
	// values.Set(_EVENT, request.Event)
	values.Set(_UPLOADED, strconv.Itoa(0))
	values.Set(_DOWNLOADED, strconv.Itoa(0))

	url.RawQuery = values.Encode()

	// Send GET request
	httpResponse, err := http.Get(url.String())

	if err != nil {
		return TrackResponse{}, fmt.Errorf("an error occurred while contacting tracker: %w", err)
	}

	response, err := parseResponse(httpResponse.Body)

	if err != nil {
		return TrackResponse{}, fmt.Errorf("an error occurred while contacting tracker: %w", err)
	}

	return response, nil
}

func parseResponse(body io.ReadCloser) (TrackResponse, error) {
	defer body.Close()

	buffer := make([]byte, common.BUFFER_SIZE)
	bytesRead, err := body.Read(buffer)

	if err != nil && err.Error() != "EOF" {
		return TrackResponse{}, err
	}

	responseBytes := buffer[:bytesRead]
	var responseDict map[string]interface{}

	err = json.Unmarshal(responseBytes, &responseDict)
	if err != nil {
		return TrackResponse{}, err
	}

	response, err := extractResponse(responseDict)
	if err != nil {
		return TrackResponse{}, err
	}

	return response, nil
}

func extractResponse(responseDict map[string]interface{}) (TrackResponse, error) {
	failureReason, isPresentFailureReason := responseDict[_FAILURE_REASON]
	interval, isPresentInterval := responseDict[_INTERVAL]
	peers, isPresentPeers := responseDict[_PEERS]

	if !isPresentFailureReason || !isPresentInterval || !isPresentPeers {
		return TrackResponse{}, errors.New("server response was not in the expected format")
	}

	failureReasonCasted, err := common.CastTo[string](failureReason)
	if err != nil {
		return TrackResponse{}, err
	}

	intervalCasted, err := common.CastTo[float64](interval)
	if err != nil {
		return TrackResponse{}, err
	}

	peersCasted, err := common.CastTo[map[string]interface{}](peers)
	if err != nil {
		return TrackResponse{}, err
	}

	peersDict, err := buildPeersDict(peersCasted)
	if err != nil {
		return TrackResponse{}, err
	}

	response := TrackResponse{
		FailureReason: failureReasonCasted,
		Interval:      int(intervalCasted),
		Peers:         peersDict,
	}

	return response, nil
}

func buildPeersDict(peers map[string]interface{}) (map[string]common.Address, error) {
	peersDict := make(map[string]common.Address)

	for id, address := range peers {
		addressStr, err := common.CastTo[string](address)
		if err != nil {
			return nil, err
		}
		splitAddress := strings.Split(addressStr, ":")

		if len(splitAddress) < 2 {
			return nil, errors.New("invalid address")
		}

		peerAddress := common.Address{
			Ip:   splitAddress[0],
			Port: splitAddress[1],
		}

		peersDict[id] = peerAddress
	}

	return peersDict, nil
}
