package server

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	userGen "github.com/reivaj05/apigateway/api/users/generated"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ServerTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	httpPort   string
	grpcPort   string
	grpcServer *grpc.Server
}

type mockServiceImpl struct {
	userGen.UserServer
}

// GetUser returns a mock user
func (service *mockServiceImpl) GetUser(
	context.Context, *userGen.Id) (*userGen.Info, error) {

	return &userGen.Info{
		Email: "email",
		Id:    "id",
		Name:  "name",
		Role:  "role",
	}, nil
}

func registerMockGRPCEndpoint(grpcServer *grpc.Server) {
	userGen.RegisterUserServer(grpcServer, &mockServiceImpl{})
}

func registerMockHTTPEndpoint(ctx context.Context, mux *runtime.ServeMux,
	endpoint string, opts []grpc.DialOption) (err error) {

	return userGen.RegisterUserHandlerFromEndpoint(
		ctx, mux, endpoint, opts)
}

func registerMockWrongHTTPEndpoint(ctx context.Context, mux *runtime.ServeMux,
	endpoint string, opts []grpc.DialOption) (err error) {

	return fmt.Errorf("Mock error")
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
	suite.httpPort = ":8000"
	suite.setupEndpoints()
	suite.deployTestServer()
}

func (suite *ServerTestSuite) setupEndpoints() {
	registeredHTTPEndpoints = []registerHTTPEndpoint{
		registerMockHTTPEndpoint,
		registerMockWrongHTTPEndpoint,
	}
	registeredGRPCEndpoints = []registerGRPCEndpoint{
		registerMockGRPCEndpoint,
	}
}

func (suite *ServerTestSuite) deployTestServer() {
	suite.grpcServer = createGRPCServer()
	go startServers(createHTTPServer(), suite.grpcServer)
	time.Sleep(time.Second)
}

func (suite *ServerTestSuite) TearDownSuite() {
	suite.grpcServer.GracefulStop()
}

func (suite *ServerTestSuite) TestHTTPRequest() {
	resp, err := http.Get("http://localhost" + suite.httpPort + "/api/v1/users/1")
	fmt.Println(resp, err)
	// TODO: Parse json and check fields
	// body := read(resp)
	// jsonBody, err := parseJson(resp)
	// suite.assertJSONFields(jsonBody)
}

// func (suite *ServerTestSuite) assertJSONFields(jsonBody something) {
// 	email, _ := jsonBody.get("email")
// 	suite.assert.Equal("email", email)
// 	id, _ := jsonBody.get("id")
// 	suite.assert.Equal("id", id)
// 	name, _ := jsonBody.get("name")
// 	suite.assert.Equal("name", name)
// 	role, _ := jsonBody.get("role")
// 	suite.assert.Equal("role", role)
// }

func (suite *ServerTestSuite) TestGRPCRequest() {
	suite.grpcPort = ":10000"
	userObj, err := suite.GetUserMock(&userGen.Id{Id: "1"})
	suite.assert.Nil(err)
	suite.assert.NotNil(userObj)
	suite.assertGRPCFields(userObj)
}

func (suite *ServerTestSuite) assertGRPCFields(userObj *userGen.Info) {
	suite.assert.Equal("email", userObj.Email)
	suite.assert.Equal("id", userObj.Id)
	suite.assert.Equal("name", userObj.Name)
	suite.assert.Equal("role", userObj.Role)
}

func (suite *ServerTestSuite) TestGRPCWrongRequest() {
	suite.grpcPort = ":7575"
	user, err := suite.GetUserMock(&userGen.Id{Id: "1"})
	suite.assert.NotNil(err)
	suite.assert.Nil(user)
}

func (suite *ServerTestSuite) GetUserMock(
	userObj *userGen.Id) (*userGen.Info, error) {

	userConn, err := grpc.Dial(suite.grpcPort, grpc.WithInsecure())
	suite.assert.Nil(err)
	suite.assert.NotNil(userConn)
	defer userConn.Close()

	userClient := userGen.NewUserClient(userConn)
	return userClient.GetUser(context.TODO(), userObj)
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
