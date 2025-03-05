package TrackerNode

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/WithSocket"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type HttpTracker struct {
	node                    *Core.BruteChord[WithSocket.SocketContact]
	logger                  common.Logger
	lock                    sync.Locker
	Ip                      string
	Port                    string
	readingRequestFunctions map[string]func(*http.Request, *common.TrackRequest) error
}

func NewHttpTracker(name string, iface string, apiPort string) *HttpTracker {
	var httpTracker HttpTracker
	httpTracker.logger = *common.NewLogger(fmt.Sprintf("HTTTracker%s.log", name))
	httpTracker.lock = &sync.Mutex{}
	ip, _ := WithSocket.GetIpFromInterface(iface)
	httpTracker.Ip = ip
	httpTracker.Port = apiPort
	httpTracker.readingRequestFunctions = map[string]func(*http.Request, *common.TrackRequest) error{
		common.InfoHash: httpTracker.handleInfoHash,
		common.PeerId:   httpTracker.handlePeerId,
		common.Ip:       httpTracker.handleIp,
		common.Port:     httpTracker.handlePort,
		common.Left:     httpTracker.handleLeft,
	}
	go receiveFromMulticast(httpTracker.Ip, httpTracker.Port)
	go httpTracker.Listen()
	httpTracker.node = WithSocket.NewNodeSocket()
	return &httpTracker
}

func (tracker *HttpTracker) Listen() {
	http.HandleFunc("/announce", tracker.handleGetPeers)
	http.HandleFunc("/nodeState", tracker.handleGetNodeState)
	address := tracker.Ip + ":" + tracker.Port
	err := http.ListenAndServe(address, nil)
	if err != nil {
		tracker.logger.WriteToFileError("Listen and Serve failed %s", err.Error())
		os.Exit(1)
	}
	tracker.logger.WriteToFileOK("Tracker listening on " + address)
}

func (tracker *HttpTracker) handleGetNodeState(w http.ResponseWriter, r *http.Request) {
	result := tracker.node.GetState()
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(result)
	if err != nil {
		tracker.logger.WriteToFileError("Failed to encode response: %s", err.Error())
	}
}

func (tracker *HttpTracker) handleGetPeers(w http.ResponseWriter, r *http.Request) {
	var request common.TrackRequest
	var response common.TrackResponse
	err := r.ParseForm()
	if err != nil {
		tracker.logger.WriteToFileError("Failed to parse URL from request: %s", err.Error())
		response.FailureReason = "Failed to parse URL from request"
		return
	}
	_, isCustomClient := r.Form[common.CustomClient] // Check if the request comes from our client.
	for attribute, function := range tracker.readingRequestFunctions {
		err = function(r, &request)
		if err != nil {
			tracker.logger.WriteToFileError("Failed to read request: %s", err.Error())
			response.FailureReason = "Error when reading the attribute " + attribute + " " + err.Error()
			tracker.sendResponse(w, response, isCustomClient)
			return
		}
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

func (tracker *HttpTracker) sendResponse(w http.ResponseWriter, response common.TrackResponse, isCustom bool) {
	var responseEncoded []byte
	var err error
	if !isCustom {
		officialResponse := common.BuildOfficialResponse(response)
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

func (tracker *HttpTracker) solve(request common.TrackRequest) (common.TrackResponse, error) {
	tracker.lock.Lock()
	defer tracker.lock.Unlock()

	var ans common.TrackResponse
	ans.Interval = 10

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
