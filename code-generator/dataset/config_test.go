package GoConfig

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type configTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *configTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
}

func (suite *configTestSuite) SetupTest() {
	viper.Reset()
}

func (suite *configTestSuite) TestInitWithWrongConfigType() {
	options := createOptions("wrongConfigType", "", "")
	err := Init(options)
	suite.assert.NotNil(err)
}

func (suite *configTestSuite) TestInitWithWrongConfigFile() {
	options := createOptions("json", "nothing", "")
	err := Init(options)
	suite.assert.NotNil(err)
}

func (suite *configTestSuite) TestInitWithGoodConfigFileButNoPath() {
	options := createOptions("json", "config", "")
	err := Init(options)
	suite.assert.NotNil(err)
}

func (suite *configTestSuite) TestInitWithCorruptedFileConfig() {
	options := createOptions("json", "corrupted", "./")
	err := Init(options)
	suite.assert.NotNil(err)
}

func (suite *configTestSuite) TestInitWithGoodConfigOptions() {
	options := createOptions("json", "config", "./")
	err := Init(options)
	suite.assert.Nil(err)
}

func createOptions(configType, configFile, configPath string) *ConfigOptions {
	return &ConfigOptions{
		ConfigType: configType,
		ConfigFile: configFile,
		ConfigPath: configPath,
	}
}

func (suite *configTestSuite) TestGetStringValue() {
	setupGetValues()
	goodKey, wrongKey := "stringKey", "wrongKey"
	value := GetConfigStringValue(goodKey)
	suite.assert.Equal("stringValue", value)
	value = GetConfigStringValue(wrongKey)
	suite.assert.Equal("", value)
}

func (suite *configTestSuite) TestGetIntValue() {
	setupGetValues()
	goodKey, wrongKey := "intKey", "wrongKey"
	value := GetConfigIntValue(goodKey)
	suite.assert.Equal(10, value)
	value = GetConfigIntValue(wrongKey)
	suite.assert.Equal(0, value)
}

func (suite *configTestSuite) TestGetFloatValue() {
	setupGetValues()
	goodKey, wrongKey := "floatKey", "wrongKey"
	value := GetConfigFloatValue(goodKey)
	suite.assert.Equal(5.5, value)
	value = GetConfigFloatValue(wrongKey)
	suite.assert.Equal(0.0, value)
}

func (suite *configTestSuite) TestGetBoolValue() {
	setupGetValues()
	goodKey, wrongKey := "booleanKey", "wrongKey"
	value := GetConfigBoolValue(goodKey)
	suite.assert.Equal(true, value)
	value = GetConfigBoolValue(wrongKey)
	suite.assert.Equal(false, value)
}

func (suite *configTestSuite) TestGetMapValue() {
	setupGetValues()
	goodKey := "mapKey"
	value := GetConfigMapValue(goodKey)
	suite.assert.Equal(float64(1), value["map1"].(float64))
	suite.assert.Equal(float64(2), value["map2"].(float64))
}

func (suite *configTestSuite) TestSetConfigValue() {
	setupGetValues()
	key := "intKey"
	oldValue := GetConfigIntValue(key)
	newValue := 100
	SetConfigValue(key, newValue)
	suite.assert.NotEqual(oldValue, GetConfigIntValue(key))
	suite.assert.Equal(newValue, GetConfigIntValue(key))
}

func (suite *configTestSuite) TestHasKey() {
	setupGetValues()
	suite.assert.True(HasKey("stringKey"))
	suite.assert.False(HasKey("wrongKey"))
}

func setupGetValues() {
	Init(createOptions("json", "config", "./"))
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(configTestSuite))
}
