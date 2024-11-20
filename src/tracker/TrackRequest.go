package tracker

type TrackRequest struct {
	InfoHash []byte // The 20 byte sha1 hash of the bencoded form of the info value from the metainfo(.torrent) file
	PeerId   string // A string of length 20 which the downloader uses as its id
	Ip       string // An optional parameter giving the IP (or dns name) which this peer is at
	Port     string // The port number this peer is listening on
	Left     int    // Number of bytes this peer has to download (Notice that if Left == 0 then the downloader becomes a seeder)
	Event    string // This is an optional key which maps to started, completed, or stopped

	// Uploaded int // Statistical property not important for functioning purposes
	// Downloaded int // Statistical property not important for functioning purposes
}
