package network

type GetBlocksMessage struct {
	From uint64
	// If To is 0 the maximum blocks will be returned.
	To uint64
}

type GetStatusMessage struct{}

type StatusMessage struct {
	// the id of the server
	ID            string
	Version       uint64
	CurrentHeight uint64
}
