package main

import (
	"errors"
	"fmt"
	"net"
	"os"

	"mockNode/api"
	"mockNode/api/pb"
	"mockNode/service"
	"mockNode/utils"

	"github.com/spacemeshos/smutil/log"
	"google.golang.org/grpc"
)

const (
	grpcDefaultAddr 	= "localhost"
	grpcDefaultPort 	= 9091
)

// parse listening address from os.Args slice, if no address was
// supplied use localhost:9091 default address
func getListenAddress(args []string) (string, error) {
	rpcListen := ""
	argsLen := len(args)
	if len(os.Args) > 2 {
		fmt.Println("Usage:", args[0], "[address:port]")
		return "", errors.New("too many arguments")
	}
	// if no arguments were given use the default address
	if argsLen == 1 {
		rpcListen = fmt.Sprintf("%v:%d", grpcDefaultAddr, grpcDefaultPort)
		return rpcListen, nil
	}

	rpcListen = os.Args[1]
	// validate address
	if !utils.ValidateFullAddress(rpcListen) {
		fmt.Println("Usage:", os.Args[0], "[address:port]")
		return "", errors.New("invalid address")
	}

	return rpcListen, nil
}

// mockNode main is the entry for creating a mock node to propagate poets proof.
// mockNode receives an optional argument in the form of "address:port",
// the full address on which the node will be listening.
func main() {
	rpcListen, err := getListenAddress(os.Args)
	if err != nil {
		fmt.Println("an has occurred error while parsing listening address", err)
		return
	}

	lis, err := net.Listen("tcp", rpcListen)
	if err != nil {
		fmt.Println("Failed listening on", rpcListen, "\nerr:", err)
		_ = fmt.Errorf("failed to Dial. err: %v", err)
		return
	}
	defer lis.Close()

	ps := &api.NodeMockServer{
		Server: grpc.NewServer(),
		Network: &service.NetworkMock{},
	}

	pb.RegisterSpacemeshServiceServer(ps.Server, ps)

	fmt.Println("grpc API listening on", rpcListen)
	if err := ps.Server.Serve(lis); err != nil {
		log.Error("failed to serve: %s", err)
	}
}