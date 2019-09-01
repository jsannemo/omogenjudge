// Package client provides an RPC client to the RunService service.
package client

import (
	"flag"

	"google.golang.org/grpc"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

var (
	serverAddr = flag.String("run_server_addr", "127.0.0.1:61811", "The runner server address to listen to in the format of host:port")
)

// Returns a new client for the RunService.
func NewClient() (runpb.RunServiceClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		return nil, err
	}
	return runpb.NewRunServiceClient(conn), nil
}
