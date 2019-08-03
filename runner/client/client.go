package client

import (
	"flag"

	"github.com/google/logger"
	"google.golang.org/grpc"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

var (
	serverAddr = flag.String("run_server_addr", "127.0.0.1:61811", "The runner server address to listen to in the format of host:port")
)

var conn *grpc.ClientConn

func getConn() *grpc.ClientConn {
	var err error
	if conn == nil {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		conn, err = grpc.Dial(*serverAddr, opts...)
		if err != nil {
			logger.Fatalf("fail to dial: %v", err)
		}
	}
	return conn
}

func NewClient() runpb.RunServiceClient {
	return runpb.NewRunServiceClient(getConn())
}
