package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type httpServerOptions struct {
	ctx     context.Context
	address string
	opts    []grpc.DialOption
}

// Serve starts grpc and http servers
func Serve(args ...string) error {
	return startServers(createHTTPServer(), createGRPCServer())
}

func createHTTPServer() *http.ServeMux {
	return addGWServerToHTTPServer(http.NewServeMux(), createGWServer())
}

func createGWServer() *runtime.ServeMux {
	return registerHTTPEndpoints(runtime.NewServeMux())
}

func registerHTTPEndpoints(gwServer *runtime.ServeMux) *runtime.ServeMux {
	config := createHTTPServerConfig()
	for _, callback := range registeredHTTPEndpoints {
		registerHTTP(config, callback, gwServer)
	}
	return gwServer
}

func createHTTPServerConfig() *httpServerOptions {
	return &httpServerOptions{
		ctx: context.Background(),
		// TODO: Get port from config file
		address: fmt.Sprintf(":%d", 10000),
		opts:    []grpc.DialOption{grpc.WithInsecure()},
	}
}

func registerHTTP(config *httpServerOptions, callback registerHTTPEndpoint,
	gwServer *runtime.ServeMux) {

	if err := callback(config.ctx, gwServer,
		config.address, config.opts); err != nil {
		// TODO: Log message and continue or finish
	}
}

func addGWServerToHTTPServer(httpServer *http.ServeMux,
	gwServer *runtime.ServeMux) *http.ServeMux {

	httpServer.Handle("/", gwServer)
	return serveSwagger(httpServer)
}

func serveSwagger(httpServer *http.ServeMux) *http.ServeMux {
	// TODO: Change to config
	swaggerDir := "./api/swagger/"
	if files, err := ioutil.ReadDir(swaggerDir); err != nil {
		// TODO: log error and continue or finish
	} else {
		for _, file := range files {
			addSwaggerEndpoint(httpServer, file, swaggerDir)
		}
	}
	return httpServer
}

func addSwaggerEndpoint(httpServer *http.ServeMux,
	file os.FileInfo, swaggerDir string) {

	data, err := ioutil.ReadFile(swaggerDir + file.Name())
	if err != nil {
		// TODO: log error and continue or finish
	}
	registerSwaggerEndpoint(httpServer, file, data)
}

func registerSwaggerEndpoint(httpServer *http.ServeMux,
	file os.FileInfo, data []byte) {

	endpoint := "/docs/" + strings.TrimSuffix(file.Name(), ".swagger.json")
	httpServer.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.Copy(w, strings.NewReader(string(data)))
	})
}

func createGRPCServer() *grpc.Server {
	return registerGRPCEndpoints(grpc.NewServer())
}

func registerGRPCEndpoints(grpcServer *grpc.Server) *grpc.Server {
	for _, callback := range registeredGRPCEndpoints {
		callback(grpcServer)
	}
	return grpcServer
}

func startServers(httpServer *http.ServeMux, grpcServer *grpc.Server) error {
	// TODO: Get port from config file
	go http.ListenAndServe(":8000", createHTTPServer())
	// TODO: Get port from config file
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", 10000))
	if err != nil {
		return err
	}
	defer conn.Close()
	return createGRPCServer().Serve(conn)
}
