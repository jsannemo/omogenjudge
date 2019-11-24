// Package client provides an RPC client to the ExecService service.
package client

import (
	"flag"

	"google.golang.org/grpc"

	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
)

var (
	serverAddr = flag.String("exec_server_addr", "127.0.0.1:61810", "The sandbox server address to connect to to in the format of host:port")
)

// NewClient creates a new client for the ExecService.
func NewClient() (execpb.ExecuteServiceClient, error) {
	var opts []grpc.DialOption
	// TODO(jsannemo): this should not use insecure credentials
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		return nil, err
	}
	return execpb.NewExecuteServiceClient(conn), nil
}
