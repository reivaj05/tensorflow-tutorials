package server

import "google.golang.org/grpc"

type registerGRPCEndpoint func(grpcServer *grpc.Server)

var registeredGRPCEndpoints = []registerGRPCEndpoint{
// users.RegisterGRPCEndpoint,
}
