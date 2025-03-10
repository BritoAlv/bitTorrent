package TrackerNode

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/WithSocket"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type HttpTracker struct {
	node   *Core.BruteChord[WithSocket.SocketContact]
	logger common.Logger
	lock   sync.Locker
	Ip     string
	Port   string
}

func NewHttpTracker(name string, iface string) *HttpTracker {
	var httpTracker HttpTracker
	httpTracker.logger = *common.NewLogger(fmt.Sprintf("HTTTracker%s.log", name))
	httpTracker.lock = &sync.Mutex{}
	ip, _ := WithSocket.GetIpFromInterface(iface)
	httpTracker.Ip = ip
	httpTracker.Port = "8080"
	go receiveFromMulticast(httpTracker.Ip, httpTracker.Port)
	go httpTracker.Listen()
	httpTracker.node = WithSocket.NewNodeSocket()
	return &httpTracker
}

func (tracker *HttpTracker) Listen() {
	http.HandleFunc("/announce", tracker.handlePeersQuery)
	http.HandleFunc("/nodeState", tracker.handleNodeState)
	address := tracker.Ip + ":" + tracker.Port
	err := http.ListenAndServe(address, nil)
	if err != nil {
		tracker.logger.WriteToFileError("Listen and Serve failed %s", err.Error())
		os.Exit(1)
	}
	tracker.logger.WriteToFileOK("Tracker listening on " + address)
}

func (tracker *HttpTracker) handleNodeState(w http.ResponseWriter, r *http.Request) {
	result := tracker.node.GetState()
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(result)
	if err != nil {
		tracker.logger.WriteToFileError("Failed to encode response: %s", err.Error())
	}
}

func (tracker *HttpTracker) handlePeersQuery(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		tracker.logger.WriteToFileError("Failed to parse URL from request: %s", err.Error())
	}
	var request common.TrackRequest
	var response common.TrackResponse

	_, isCustomClient := r.Form[common.CustomClient]

	if _, exist := r.Form[common.InfoHash]; !exist {
		message := "InfoHash not found in request"
		tracker.logger.WriteToFileError(message)
		response.FailureReason = message
		tracker.sendResponse(w, response, isCustomClient)
		return
	}
	request.InfoHash, err = common.DecodeStrByt(r.Form[common.InfoHash][0])
	if err != nil {
		tracker.logger.WriteToFileError("Failed to decode the InfoHash to bytes: %s", err.Error())
		response.FailureReason = err.Error()
		tracker.sendResponse(w, response, isCustomClient)
		return
	}

	if _, exist := r.Form[common.PeerId]; !exist {
		message := "PeerId not found in request"
		tracker.logger.WriteToFileError(message)
		response.FailureReason = message
		tracker.sendResponse(w, response, isCustomClient)
		return
	}
	request.PeerId, err = url.QueryUnescape(r.Form[common.PeerId][0])
	if err != nil {
		tracker.logger.WriteToFileError("Failed to unescape the PeerId : %s", err.Error())
		response.FailureReason = err.Error()
		tracker.sendResponse(w, response, isCustomClient)
		return
	}

	if _, exist := r.Form[common.Ip]; !exist {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			tracker.logger.WriteToFileError("Failed to split host and Port from remote address: %s", err.Error())
			response.FailureReason = err.Error()
			tracker.sendResponse(w, response, isCustomClient)
			return
		}
		request.Ip = host
		tracker.logger.WriteToFileOK("Ip not found in request, using remote address %s", request.Ip)
	} else {
		request.Ip = r.Form[common.Ip][0]
	}

	if _, exist := r.Form[common.Port]; !exist {
		message := "Port not found in request"
		tracker.logger.WriteToFileError(message)
		response.FailureReason = message
		tracker.sendResponse(w, response, isCustomClient)
		return
	}
	request.Port = r.Form[common.Port][0]

	if _, exist := r.Form[common.Left]; !exist {
		message := "Left not found in request"
		tracker.logger.WriteToFileError(message)
		response.FailureReason = message
		tracker.sendResponse(w, response, isCustomClient)
		return
	}
	request.Left, err = strconv.Atoi(r.Form[common.Left][0])
	if err != nil {
		tracker.logger.WriteToFileError("Failed to convert left to int: %s", err.Error())
		response.FailureReason = err.Error()
		tracker.sendResponse(w, response, isCustomClient)
		return
	}
	err = common.ValidateRequest(request)
	if err != nil {
		tracker.logger.WriteToFileError("Failed to validate request: %s", err.Error())
		response.FailureReason = err.Error()
		tracker.sendResponse(w, response, isCustomClient)
		return
	}
	tracker.logger.WriteToFileOK("Received request was decoded and its valid: %v", request)
	response, err = tracker.solve(request)
	if err != nil {
		tracker.logger.WriteToFileError("Failed to solve request: %s", err.Error())
		response.FailureReason = err.Error()
	}
	tracker.logger.WriteToFileOK("Will send this response: %v", response)
	tracker.sendResponse(w, response, isCustomClient)
}

func buildOfficialResponse(response common.TrackResponse) common.OfficialTrackResponse {
	var officialResponse common.OfficialTrackResponse
	officialResponse.FailureReason = response.FailureReason
	officialResponse.Interval = response.Interval
	officialResponse.Peers = make([]common.OfficialTrackAddress, 0)
	for _, address := range response.Peers {
		port, _ := strconv.Atoi(address.Port)
		officialResponse.Peers = append(officialResponse.Peers, common.OfficialTrackAddress{
			Ip:   address.Ip,
			Port: port,
		})
	}
	return officialResponse
}

func (tracker *HttpTracker) sendResponse(w http.ResponseWriter, response common.TrackResponse, isCustom bool) {
	var responseEncoded []byte
	var err error
	if !isCustom {
		officialResponse := buildOfficialResponse(response)
		responseEncoded, err = common.EncodeOfficialResponse(officialResponse)
	} else {
		responseEncoded, err = common.EncodeResponse(response)
	}
	if err != nil {
		tracker.logger.WriteToFileError("Failed to encode response: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseEncoded)
	if err != nil {
		tracker.logger.WriteToFileError("Failed to write response: %s", err.Error())
	}
	tracker.logger.WriteToFileOK("Response sent: %s", responseEncoded)
}

/*
Solve is where the tracker decides what to do with the client query for peers.
*/
func (tracker *HttpTracker) solve(request common.TrackRequest) (common.TrackResponse, error) {
	tracker.lock.Lock()
	defer tracker.lock.Unlock()

	var ans common.TrackResponse
	ans.FailureReason = ""
	ans.Interval = 100000

	infoHashToChordKey := tracker.InfoHashToChordKey(request.InfoHash)
	fmt.Println("InfoHashToChordKey", infoHashToChordKey)
	_, exist := tracker.node.Get(infoHashToChordKey)
	if !exist {
		tracker.logger.WriteToFileOK("New entry for info hash %v", request.InfoHash)
		tracker.node.Put(infoHashToChordKey, EncodePeerList(make(map[string]common.Address)))
	}

	valueInfoHash, _ := tracker.node.Get(infoHashToChordKey)
	peersInfoHash := DecodePeerList(valueInfoHash)

	if _, exist := peersInfoHash[request.PeerId]; !exist {
		tracker.logger.WriteToFileOK("New entry for peer id %v", request.PeerId)
		peersInfoHash[request.PeerId] = common.Address{
			Ip:   request.Ip,
			Port: request.Port,
		}
		tracker.node.Put(infoHashToChordKey, EncodePeerList(peersInfoHash))
	}
	ans.Peers = make(map[string]common.Address)
	for id, address := range peersInfoHash {
		if id != request.PeerId {
			ans.Peers[id] = address
		}
	}
	return ans, nil
}
