package generator

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/reivaj05/GoConfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GeneratorTestSuite struct {
	suite.Suite
	assert          *assert.Assertions
	mockServiceName string
	path            string
}

func (suite *GeneratorTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
	suite.mockServiceName = "mockService"
	suite.path = joinPath()

	GoConfig.Init(&GoConfig.ConfigOptions{
		ConfigType: "json",
		ConfigFile: "config",
		ConfigPath: "..",
	})
	GoConfig.SetConfigValue("templatesPath", "templates/")
	GoConfig.SetConfigValue("protosPath", "../protos/api/")
}

func (suite *GeneratorTestSuite) TearDownSuite() {
	rollback(suite.mockServiceName)
}

func (suite *GeneratorTestSuite) TestJoinPath() {
	suite.assert.True(strings.HasPrefix(suite.path, os.Getenv("GOPATH")))
	suite.assert.True(strings.HasSuffix(suite.path, "/src/github.com/reivaj05/apigateway"))
}

// TODO: Add unsuccessful tests
func (suite *GeneratorTestSuite) TestGenerateAPIFileSuccessful() {
	suite.assert.Nil(generateAPIFile(suite.path, suite.mockServiceName))
	_, err := os.Stat(fmt.Sprintf("../api/%s/%s.go",
		suite.mockServiceName, suite.mockServiceName))
	suite.assert.False(os.IsNotExist(err))
}

func (suite *GeneratorTestSuite) TestGenerateServiceFileSuccessful() {
	suite.assert.Nil(generateServiceFile(suite.path, suite.mockServiceName))
	_, err := os.Stat(fmt.Sprintf("../services/%s/%s.go",
		suite.mockServiceName, suite.mockServiceName))
	suite.assert.False(os.IsNotExist(err))
}

func (suite *GeneratorTestSuite) TestGenerateProtoFilesSuccessful() {
	suite.assert.Nil(generateProtoFiles(suite.path, suite.mockServiceName))
	_, err := os.Stat(fmt.Sprintf("../protos/api/%s.proto", suite.mockServiceName))
	suite.assert.False(os.IsNotExist(err))
	_, err = os.Stat(fmt.Sprintf("../protos/services/%s.proto", suite.mockServiceName))
	suite.assert.False(os.IsNotExist(err))
}

func (suite *GeneratorTestSuite) TestGetServicesNamesPopulatedList() {
	generateAPIFile(suite.path, suite.mockServiceName)
	services := getServicesNames()
	suite.assert.Equal(1, len(services))
}

func (suite *GeneratorTestSuite) TestUpdateServerFilesWrongLocation() {
	err := updateServerFiles(&EndpointsData{Services: getServicesNames()})
	suite.assert.NotNil(err)
}

func (suite *GeneratorTestSuite) TestUpdateServerFilesGoodLocation() {
	serverPath := GoConfig.GetConfigStringValue("serverPath")
	GoConfig.SetConfigValue("serverPath", "../server/")
	err := updateServerFiles(&EndpointsData{Services: getServicesNames()})
	suite.assert.Nil(err)
	GoConfig.SetConfigValue("serverPath", serverPath)
}

func (suite *GeneratorTestSuite) TestGenerateWithoutArgs() {
	err := Generate()
	suite.assert.NotNil(err)
}

func (suite *GeneratorTestSuite) TestGenerateWithArgs() {
	serverPath := GoConfig.GetConfigStringValue("serverPath")
	GoConfig.SetConfigValue("serverPath", "../server/")
	GoConfig.SetConfigValue("protoGenPath", "../proto-gen.sh")
	err := Generate(suite.mockServiceName)
	suite.assert.Nil(err)
	GoConfig.SetConfigValue("serverPath", serverPath)
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}
