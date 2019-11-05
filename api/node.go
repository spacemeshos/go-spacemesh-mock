package api

type NetworkAPI interface {
	Broadcast(channel string, data []byte) error
	GetProof(roundNum int) ([]byte, error)
}

