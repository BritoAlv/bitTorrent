package TrackerNode

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

func (tracker *HttpTracker) InfoHashToChordKey(infoHash [20]byte) Core.ChordHash {
	sum := 0
	for i := 0; i < 20; i++ {
		sum += int(infoHash[i])
	}
	return Core.ChordHash(sum % (1 << Core.NumberBits))
}

func EncodePeerList(peers map[string]common.Address) []byte {
	fmt.Printf("Passed to encode this %v\n", peers)
	if peers == nil {
		panic("Passed Peers is nil")
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(peers)
	if err != nil {
		panic(err)
	}
	bytesEncoded := buf.Bytes()
	fmt.Printf("Encoded this \n%v\n", bytesEncoded)
	return bytesEncoded
}

func DecodePeerList(data []byte) map[string]common.Address {
	var peers map[string]common.Address
	fmt.Printf("Received this to decode \n%v\n", data)
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&peers)
	if err != nil {
		panic(err)
	}
	return peers
}

func (tracker *HttpTracker) handleInfoHash(r *http.Request, request *common.TrackRequest) error {
	if _, exist := r.Form[common.InfoHash]; !exist {
		message := "InfoHash not found in request"
		tracker.logger.WriteToFileError(message)
		return fmt.Errorf(message)
	}
	infoHash, err := common.DecodeStrByt(r.Form[common.InfoHash][0])
	if err != nil {
		return err
	}
	request.InfoHash = infoHash
	return nil
}

func (tracker *HttpTracker) handlePeerId(r *http.Request, request *common.TrackRequest) error {
	if _, exist := r.Form[common.PeerId]; !exist {
		message := "PeerId not found in request"
		tracker.logger.WriteToFileError(message)
		return fmt.Errorf(message)
	}
	peerId, err := url.QueryUnescape(r.Form[common.PeerId][0])
	if err != nil {
		tracker.logger.WriteToFileError("Failed to unescape the PeerId : %s", err.Error())
		return err
	}
	request.PeerId = peerId
	if len(request.PeerId) != 20 {
		return fmt.Errorf("invalid peer id field it should have 20 bytes %s", request.PeerId)
	}
	return nil
}

func (tracker *HttpTracker) handleIp(r *http.Request, request *common.TrackRequest) error {
	if _, exist := r.Form[common.Ip]; !exist {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			tracker.logger.WriteToFileError("Failed to split host and Port from remote address: %s", err.Error())
			return err
		}
		request.Ip = host
		tracker.logger.WriteToFileOK("Ip not found in request, using remote address %s", request.Ip)
	} else {
		request.Ip = r.Form[common.Ip][0]
	}
	if net.ParseIP(request.Ip) == nil {
		return fmt.Errorf("invalid ip field %s", request.Ip)
	}
	return nil
}

func (tracker *HttpTracker) handlePort(r *http.Request, request *common.TrackRequest) error {
	if _, exist := r.Form[common.Port]; !exist {
		message := "port not found in request"
		return fmt.Errorf(message)
	}
	request.Port = r.Form[common.Port][0]
	port, err := strconv.Atoi(request.Port)
	if err != nil {
		return fmt.Errorf("invalid port field %w", err)
	}
	if port < 0 || port > 65535 {
		return fmt.Errorf("invalid port field %d", port)
	}
	return nil
}
func (tracker *HttpTracker) handleLeft(r *http.Request, request *common.TrackRequest) error {
	if _, exist := r.Form[common.Left]; !exist {
		message := "left not found in request"
		tracker.logger.WriteToFileError(message)
		return fmt.Errorf(message)
	}
	left, err := strconv.Atoi(r.Form[common.Left][0])
	if err != nil {
		tracker.logger.WriteToFileError("Failed to convert left to int: %s", err.Error())
		return fmt.Errorf(err.Error())
	}
	request.Left = left
	if request.Left < 0 {
		return fmt.Errorf("invalid left field")
	}
	return nil
}
