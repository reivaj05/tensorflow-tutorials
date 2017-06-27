package GoLogger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type loggerTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *loggerTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
}

func (suite *loggerTestSuite) TestSuccessfulInit() {
	options := createOptions("log/", "logger.test.json", WARNING)
	suite.assert.Nil(Init(options))
	suite.assert.Equal(WARNING, GetLogLevel())
}

func (suite *loggerTestSuite) TestSuccessfulInitWithWrongLogLevel() {
	options := createOptions("log/", "logger.test.json", -1)
	suite.assert.Nil(Init(options))
	suite.assert.Equal(DEBUG, GetLogLevel())
}

func (suite *loggerTestSuite) TestUnsuccessfulInitWithoutOutputFile() {
	options := createOptions("", "", DEBUG)
	suite.assert.NotNil(Init(options))
}

func (suite *loggerTestSuite) TestSetGetLogLevel() {
	SetLogLevel(PANIC)
	suite.assert.Equal(PANIC, GetLogLevel())
	SetLogLevel(DEBUG)
	suite.assert.Equal(DEBUG, GetLogLevel())
}

func (suite *loggerTestSuite) TestLogAllLevelsWithoutFields() {
	Init(createOptions("log/", "logger.test.json", DEBUG))
	msg := "Test message"
	LogError(msg, nil)
	LogWarning(msg, nil)
	LogInfo(msg, nil)
	LogDebug(msg, nil)
}

func (suite *loggerTestSuite) TestLogAllLevelsWithFields() {
	Init(createOptions("log/", "logger.test.json", DEBUG))
	msg := "Test message"
	fields := createFields()
	LogError(msg, fields)
	LogWarning(msg, fields)
	LogInfo(msg, fields)
	LogDebug(msg, fields)
}

func createFields() map[string]interface{} {
	return map[string]interface{}{
		"string":  "string",
		"int":     1,
		"float":   2.5,
		"boolean": true,
	}
}

func createOptions(path, file string, level int) *LoggerOptions {
	return &LoggerOptions{
		LogLevel:   level,
		Path:       path,
		OutputFile: file,
	}
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(loggerTestSuite))
}
