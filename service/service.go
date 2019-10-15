package service

import (
	"errors"
	"fmt"

	"github.com/spacemeshos/smutil/log"
)

// networkMocks is a list with all history of rounds broadcasted data.
// starting the list at index 1 so every round will be appended
// in the list with its' index (no round 0).
var networkMocks = make([]NetworkMock, 1)
var roundNumber = 0

// NetworkMock is a struct holding a rounds data, it includes the
// round number, it's binary data and a bool pointing whether there
// was an error at the current round
type NetworkMock struct {
	roundNumber	 int
	broadCastErr bool
	broadcasted  []byte
}

// Broadcast saves the rounds' data, NO ERROR return value may be returned
// in this format.
func (nm *NetworkMock) Broadcast(chanel string, payload []byte) error {
	log.Info("Broadcasting poet proof")
	// increment round number, starting from round 1
	// and not 0
	roundNumber += 1
	// saving broadcast data
	nm.broadcasted = payload
	nm.roundNumber = roundNumber
	networkMocks = append(networkMocks, *nm)
	if nm.broadCastErr {
		return errors.New("error during broadcast")
	}

	return nil
}

// GetProof receives an integer representing the rounds' number and returning the
// corresponding proof value of the same round
func (nm *NetworkMock) GetProof(inRoundNum int) ([]byte, error) {
	log.Info("Getting the proof of round %v", inRoundNum)
	// validate index doesn't exceeds list length
	if inRoundNum > len(networkMocks) {
		errMsg := fmt.Sprintf(
			"index out of range, index number: %v, list len: %v", inRoundNum, len(networkMocks))
		return nil, errors.New(errMsg)
	}

	if inRoundNum < 1 {
		errMsg := fmt.Sprintf("Minimal round number is 1")
		return nil, errors.New(errMsg)
	}

	return networkMocks[inRoundNum].broadcasted, nil
}
