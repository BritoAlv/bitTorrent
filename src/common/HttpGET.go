package common

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/zeebo/bencode"
)

const InfoHash = "info_hash"
const PeerId = "peer_id"
const Ip = "ip"
const Port = "port"
const Left = "left"

func DecodeStrByt(s string) ([20]byte, error) {
	s1, err := url.QueryUnescape(s)
	if err != nil {
		return [20]byte{}, err
	}
	var result [20]byte
	copy(result[:], s1)
	return result, nil
}

// EncodeResponse /*
func EncodeResponse(response TrackResponse) ([]byte, error) {
	result, err := bencode.EncodeBytes(response)
	if err != nil {
		return []byte{}, err
	}
	return result, nil
}

// BuildHttpUrl /*
func BuildHttpUrl(trackerUrl string, request TrackRequest) (string, error) {
	trackerL, err := url.Parse(trackerUrl)
	if err != nil {
		fmt.Println("an error occurred while contacting tracker: %w", err)
		return "", err
	}
	values := trackerL.Query()
	values.Set(InfoHash, url.QueryEscape(string(request.InfoHash[:])))
	values.Set(PeerId, url.QueryEscape(request.PeerId))
	values.Set(Ip, request.Ip)
	values.Set(Port, request.Port)
	values.Set(Left, strconv.Itoa(request.Left))
	trackerL.RawQuery = values.Encode()
	return trackerL.String(), nil
}

// DecodeTrackerResponse /*
func DecodeTrackerResponse(bytes []byte) (TrackResponse, error) {
	var response TrackResponse
	err := bencode.DecodeBytes(bytes, &response)
	if err != nil {
		return TrackResponse{}, err
	}
	return response, nil
}

// ValidateRequest  /*
func ValidateRequest(request TrackRequest) error {
	if request.Left < 0 {
		return fmt.Errorf("invalid left field")
	}
	port, err := strconv.Atoi(request.Port)
	if err != nil {
		return fmt.Errorf("invalid port field %w", err)
	}
	if port < 0 || port > 65535 {
		return fmt.Errorf("invalid port field %d", port)
	}
	if net.ParseIP(request.Ip) == nil {
		return fmt.Errorf("invalid ip field %s", request.Ip)
	}
	if len(request.PeerId) != 20 {
		return fmt.Errorf("invalid peer id field it should have 20 bytes %s", request.PeerId)
	}
	return nil
}