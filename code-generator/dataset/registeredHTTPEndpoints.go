package server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type registerHTTPEndpoint func(ctx context.Context, mux *runtime.ServeMux,
	endpoint string, opts []grpc.DialOption) (err error)

var registeredHTTPEndpoints = []registerHTTPEndpoint{
// users.RegisterHTTPEndpoint,
}
