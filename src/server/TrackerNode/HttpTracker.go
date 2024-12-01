package TrackerNode

import (
	"bittorrent/common"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type HttpTracker struct {
	peers    map[[20]byte]map[string]common.Address
	logger   common.Logger
	lock     sync.Locker
	Location string
}

func NewHttpTracker(url string, name string) TrackerNode {
	return HttpTracker{
		peers:    make(map[[20]byte]map[string]common.Address),
		logger:   *common.NewLogger(fmt.Sprintf("HTTTracker%s.log", name)),
		lock:     &sync.Mutex{},
		Location: url,
	}
}

func (tracker HttpTracker) SaveTorrent() string {
	return "http://localhost:" + tracker.Location + "/announce"
}

func (tracker HttpTracker) Listen() error {
	http.HandleFunc("/announce", tracker.handlePeersQuery)
	err := http.ListenAndServe(":"+tracker.Location, nil)
	if err != nil {
		tracker.logger.WriteToFileError(fmt.Sprintf("Listen and Serve failed %s", err.Error()))
		return err
	}
	tracker.logger.WriteToFileOK("Tracker listening on " + tracker.Location)
	return nil
}

/*
Tracker have to handle the query from the client.
*/
func (tracker HttpTracker) handlePeersQuery(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		tracker.logger.WriteToFileError(fmt.Sprintf("Failed to parse URL from request: %s", err.Error()))
	}
	var request common.TrackRequest
	var response common.TrackResponse

	if _, exist := r.Form[common.InfoHash]; !exist {
		message := "InfoHash not found in request"
		tracker.logger.WriteToFileError(message)
		response.FailureReason = message
		tracker.sendResponse(w, response)
		return
	}
	request.InfoHash, err = common.DecodeStrByt(r.Form[common.InfoHash][0])
	if err != nil {
		tracker.logger.WriteToFileError(fmt.Sprintf("Failed to decode the InfoHash to bytes: %s", err.Error()))
		response.FailureReason = err.Error()
		tracker.sendResponse(w, response)
		return
	}

	if _, exist := r.Form[common.PeerId]; !exist {
		message := "PeerId not found in request"
		tracker.logger.WriteToFileError(message)
		response.FailureReason = message
		tracker.sendResponse(w, response)
		return
	}
	request.PeerId, err = url.QueryUnescape(r.Form[common.PeerId][0])
	if err != nil {
		tracker.logger.WriteToFileError(fmt.Sprintf("Failed to unescape the PeerId : %s", err.Error()))
		response.FailureReason = err.Error()
		tracker.sendResponse(w, response)
		return
	}

	if _, exist := r.Form[common.Ip]; !exist {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			tracker.logger.WriteToFileError(fmt.Sprintf("Failed to split host and port from remote address: %s", err.Error()))
			response.FailureReason = err.Error()
			tracker.sendResponse(w, response)
			return
		}
		request.Ip = host
		tracker.logger.WriteToFileOK(fmt.Sprintf("Ip not found in request, using remote address %s", request.Ip))
	} else {
		request.Ip = r.Form[common.Ip][0]
	}

	if _, exist := r.Form[common.Port]; !exist {
		message := "Port not found in request"
		tracker.logger.WriteToFileError(message)
		response.FailureReason = message
		tracker.sendResponse(w, response)
		return
	}
	request.Port = r.Form[common.Port][0]

	if _, exist := r.Form[common.Left]; !exist {
		message := "Left not found in request"
		tracker.logger.WriteToFileError(message)
		response.FailureReason = message
		tracker.sendResponse(w, response)
		return
	}
	request.Left, err = strconv.Atoi(r.Form[common.Left][0])
	if err != nil {
		tracker.logger.WriteToFileError(fmt.Sprintf("Failed to convert left to int: %s", err.Error()))
		response.FailureReason = err.Error()
		tracker.sendResponse(w, response)
		return
	}
	err = common.ValidateRequest(request)
	if err != nil {
		tracker.logger.WriteToFileError(fmt.Sprintf("Failed to validate request: %s", err.Error()))
		response.FailureReason = err.Error()
		tracker.sendResponse(w, response)
		return
	} 		
	tracker.logger.WriteToFileOK(fmt.Sprintf("Received request was decoded and its valid: %v", request))
	response, err = tracker.solve(request)
	if err != nil {
		tracker.logger.WriteToFileError(fmt.Sprintf("Failed to solve request: %s", err.Error()))
		response.FailureReason = err.Error()
	}
	
	tracker.logger.WriteToFileOK(fmt.Sprintf("Will send this response: %v", response))
	tracker.sendResponse(w, response)
}

func (tracker HttpTracker) sendResponse(w http.ResponseWriter, response common.TrackResponse) {
	responseEncoded, err := common.EncodeResponse(response)
	if err != nil {
		tracker.logger.WriteToFileError(fmt.Sprintf("Failed to encode response: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseEncoded)
	if err != nil {
		tracker.logger.WriteToFileError(fmt.Sprintf("Failed to write response: %s", err.Error()))
	}
	tracker.logger.WriteToFileOK(fmt.Sprintf("Response sent: %s", responseEncoded))
}

/*
Solve is where the tracker decides what to do with the client query for peers.
*/
func (tracker HttpTracker) solve(request common.TrackRequest) (common.TrackResponse, error) {
	tracker.lock.Lock()
	defer tracker.lock.Unlock()

	var ans common.TrackResponse
	ans.FailureReason = ""
	ans.Interval = 10000

	if _, exist := tracker.peers[request.InfoHash]; !exist {
		tracker.logger.WriteToFileOK(fmt.Sprintf("New entry for info hash %v", request.InfoHash))
		tracker.peers[request.InfoHash] = make(map[string]common.Address)
	}

	if _, exist := tracker.peers[request.InfoHash][request.PeerId]; !exist {
		tracker.logger.WriteToFileOK(fmt.Sprintf("New entry for peer id %v", request.PeerId))
		tracker.peers[request.InfoHash][request.PeerId] = common.Address{
			Ip:   request.Ip,
			Port: request.Port,
		}
	}
	ans.Peers = make(map[string]common.Address)
	for id, address := range tracker.peers[request.InfoHash] {
		if id != request.PeerId {
			ans.Peers[id] = address
		}
	}
	return ans, nil
}