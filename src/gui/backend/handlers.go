package backend

import (
	"bittorrent/common"
	"encoding/json"
	"io"
	"net/http"
)

type DownloadHandler struct{}
type UpdateHandler struct{}
type KillHandler struct{}

func (handler *DownloadHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := getBody(request.Body)
	if err != nil {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: "Internal server error " + err.Error(),
		})
		return
	}

	var downloadRequest DownloadRequest
	err = json.Unmarshal(body, &downloadRequest)
	if err != nil {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: "Internal server error " + err.Error(),
		})
		return
	}

	respond(responseWriter, BooleanResponse{
		Successful:   true,
		ErrorMessage: "",
	})
}

func (handler *UpdateHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	queryMap := request.URL.Query()
	id := queryMap.Get("id")

	if id == "" {
		respond(responseWriter, UpdateResponse{
			BooleanResponse: BooleanResponse{
				Successful:   false,
				ErrorMessage: "Expecting 'id' query",
			},
			Progress: 0,
			Peers:    0,
		})
		return
	}

	respond(responseWriter, UpdateResponse{
		BooleanResponse: BooleanResponse{
			Successful:   true,
			ErrorMessage: "",
		},
		Progress: 0.5,
		Peers:    10,
	})
}

func (handler *KillHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	queryMap := request.URL.Query()
	id := queryMap.Get("id")

	if id == "" {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: "Expecting 'id' query",
		})
		return
	}

	respond(responseWriter, BooleanResponse{
		Successful:   true,
		ErrorMessage: "Killed" + id,
	})
}

func getBody(body io.ReadCloser) ([]byte, error) {
	buffer := make([]byte, 2048)
	totalRead, err := body.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	return buffer[:totalRead], nil
}

func respond(responseWriter io.Writer, response interface{}) {
	responseBytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: "Internal server error" + err.Error(),
		})
		return
	}

	err = common.ReliableWrite(responseWriter, responseBytes)
	if err != nil {
		respond(responseWriter, BooleanResponse{
			Successful:   false,
			ErrorMessage: "Internal server error" + err.Error(),
		})
		return
	}
}

// func getFile(path string) (*os.File, error) {
// 	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, common.ALL_RW_PERMISSION)
// 	if err != nil {
// 		return nil, errors.New("error accessing the file: " + err.Error())
// 	}

// 	return file, nil
// }
