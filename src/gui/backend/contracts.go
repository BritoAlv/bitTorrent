package backend

type DownloadRequest struct {
	Id              string
	TorrentPath     string
	DownloadPath    string
	IpAddress       string
	EncryptionLevel bool
}

type BooleanResponse struct {
	Successful   bool
	ErrorMessage string
}

type UpdateResponse struct {
	BooleanResponse
	Progress float32
	Peers    int
}
