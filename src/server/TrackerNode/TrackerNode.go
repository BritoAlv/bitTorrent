package TrackerNode

type TrackerNode interface {
	// SaveTorrent : string
	// When a Torrent is being created and wants to add this tracker to its list,
	// it should get how the tracker announces itself.
	SaveTorrent() string
	// Listen
	// Starts the Tracker.
	Listen() error
}
