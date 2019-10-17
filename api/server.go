package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spacemeshos/go-spacemesh-mock/api/nmpb"
	"github.com/spacemeshos/smutil/log"

	"google.golang.org/grpc"
)

const PoetProofProtocol	= "PoetProof"

type NodeMockServer struct {
	Server  *grpc.Server
	Network NetworkAPI
}

// BroadcastPoet method broadcasts poets proofs, it uses the network
// API broadcast method to do so
func (ns NodeMockServer) BroadcastPoet(ctx context.Context, in *nmpb.BinaryMessage) (*nmpb.SimpleMessage, error) {
	log.Info("GRPC Broadcast PoET msg")
	// call NetworkMock service broadcast method
	err := ns.Network.Broadcast(PoetProofProtocol, in.Data)
	if err != nil {
		log.Error("error in BroadcastPoet: %v", err)
		return &nmpb.SimpleMessage{Value: err.Error()}, err
	}
	log.Debug("PoET message broadcast succeeded")
	return &nmpb.SimpleMessage{Value: "ok"}, nil
}

// GetProof gets a message with an integer value (represented by a string)
// and using network api to return the corresponding proof value of the same round
func (ns NodeMockServer) GetProof(ctx context.Context, roundNumMsg *nmpb.SimpleMessage) (*nmpb.BinaryMessage, error) {
	roundNum, err := strconv.Atoi(roundNumMsg.Value)
	if err != nil {
		fmt.Println("Error converting round number from string to int:", err)
		return nil, err
	}
	// call service broadcast method
	data, err := ns.Network.GetProof(roundNum)
	if err != nil {
		log.Error("error in GetProof: %v", err)
		return nil, err
	}
	log.Debug("Proof request succeeded")
	return &nmpb.BinaryMessage{Data: data}, nil
}
