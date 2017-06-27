package GoJSON

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JSONWrapperTestSuite struct {
	suite.Suite
	Assert        *assert.Assertions
	TestErrorJSON string
}

func (suite *JSONWrapperTestSuite) SetupSuite() {
	suite.Assert = assert.New(suite.T())

	suite.TestErrorJSON = `{"TestError":}`
}

func (suite *JSONWrapperTestSuite) TestNew() {
	json, err := New(suite.TestErrorJSON)
	suite.Assert.Nil(json)
	suite.Assert.NotNil(err)

	json, err = New("{}")
	suite.Assert.NotNil(json)
	suite.Assert.Nil(err)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetStringFromPath() {
	json, _ := New("{}")
	json.SetValueAtPath("testString", "String")
	data, _ := json.GetStringFromPath("testString")
	suite.Assert.Equal("String", data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetString() {
	json, _ := New("{}")
	json.SetValueAtPath("testString", "String")
	jsonObject := json.GetJSONObjectFromPath("testString")
	data, _ := jsonObject.GetString()
	suite.Assert.Equal("String", data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetIntFromPath() {
	json, _ := New("{}")
	json.SetValueAtPath("testInt", 10)
	data, _ := json.GetIntFromPath("testInt")
	suite.Assert.Equal(10, data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetInt() {
	json, _ := New("{}")
	json.SetValueAtPath("testInt", 10)
	jsonObject := json.GetJSONObjectFromPath("testInt")
	data, _ := jsonObject.GetInt()
	suite.Assert.Equal(10, data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetInt64FromPath() {
	json, _ := New("{}")
	json.SetValueAtPath("testInt64", int64(25))
	data, _ := json.GetInt64FromPath("testInt64")
	suite.Assert.Equal(int64(25), data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetInt64() {
	json, _ := New("{}")
	json.SetValueAtPath("testInt64", int64(25))
	jsonObject := json.GetJSONObjectFromPath("testInt64")
	data, _ := jsonObject.GetInt64()
	suite.Assert.Equal(int64(25), data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetFloatFromPath() {
	json, _ := New("{}")
	json.SetValueAtPath("testFloat", 10.5)
	data, _ := json.GetFloatFromPath("testFloat")
	suite.Assert.Equal(10.5, data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetFloat() {
	json, _ := New("{}")
	json.SetValueAtPath("testFloat", 10.5)
	jsonObject := json.GetJSONObjectFromPath("testFloat")
	data, _ := jsonObject.GetFloat()
	suite.Assert.Equal(10.5, data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetBoolFromPath() {
	json, _ := New("{}")
	json.SetValueAtPath("testBool", true)
	data, _ := json.GetBoolFromPath("testBool")
	suite.Assert.Equal(true, data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetBool() {
	json, _ := New("{}")
	json.SetValueAtPath("testBool", true)
	jsonObject := json.GetJSONObjectFromPath("testBool")
	data, _ := jsonObject.GetBool()
	suite.Assert.Equal(true, data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetJSONObjectFromPath() {
	json, _ := New("{}")
	json.SetValueAtPath("testBool", true)
	obj := json.GetJSONObjectFromPath("testBool")
	suite.Assert.NotNil(obj)

	obj = json.GetJSONObjectFromPath("fake")
	suite.Assert.Nil(obj)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetJSON() {
	json, _ := New("{}")
	obj := json.GetJSONObject()
	suite.Assert.NotNil(obj)
	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestCopyJSONObjectFromPath() {
	json, _ := New("{}")
	json.SetValueAtPath("testBool", true)
	obj := json.CopyJSONObjectFromPath("testBool")
	suite.Assert.NotNil(obj)

	obj = json.CopyJSONObjectFromPath("fake")
	suite.Assert.Nil(obj)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetArrayFromPath() {
	json, _ := New(`{"testArray": [1, 2, 3]}`)
	data := json.GetArrayFromPath("testArray")
	for index, item := range data {
		value, _ := item.GetFloat()
		suite.Assert.Equal(float64(index+1), value)
	}

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestGetArray() {
	json, _ := New(`{"testArray": [1, 2, 3]}`)
	jsonObject := json.GetJSONObjectFromPath("testArray")
	data := jsonObject.GetArray()
	for index, item := range data {
		value, _ := item.GetFloat()
		suite.Assert.Equal(float64(index+1), value)
	}

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestArrayAppendCopy() {
	json, _ := New(`{"testArray": [1, 2, 3], "testInt": 1}`)
	jsonObject := json.GetJSONObjectFromPath("testArray")
	jsonInt := json.GetJSONObjectFromPath("testInt")
	suite.Assert.Nil(jsonObject.ArrayAppendCopy(jsonInt))

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestSetValueAtPath() {
	json, _ := New("{}")
	path := "testPath"
	value := "testValue"

	json.SetValueAtPath(path, value)
	data, _ := json.GetStringFromPath(path)
	suite.Assert.Equal(value, data)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestSetObjectAtPath() {
	json, _ := New(`{"testString": "String"}`)
	pathTo := "testPath"
	pathFrom := "testString"

	json.SetObjectAtPath(pathTo, json.GetJSONObjectFromPath(pathFrom))

	expected, _ := json.GetStringFromPath(pathFrom)
	actual, _ := json.GetStringFromPath(pathTo)
	suite.Assert.Equal(expected, actual)

	json.FreeJSON()
}

func (suite *JSONWrapperTestSuite) TestToString() {
	mockJSON := `{"testCount":5}`
	json, _ := New(mockJSON)
	defer json.FreeJSON()
	suite.Assert.Equal(mockJSON, json.ToString())
}

func (suite *JSONWrapperTestSuite) TestHasPath() {
	mockJSON := `{"testCount":5}`
	json, _ := New(mockJSON)
	defer json.FreeJSON()
	suite.Assert.True(json.HasPath("testCount"))
	suite.Assert.False(json.HasPath("something"))
}

func (suite *JSONWrapperTestSuite) TestCreateJSONArrayAtPathWithArray() {
	json, _ := New(`{"testArray": [1, 2, 3]}`)
	err := json.CreateJSONArrayAtPathWithArray(
		"newArray", json.GetArrayFromPath("testArray"))
	suite.Assert.Nil(err)
	suite.Assert.True(json.HasPath("newArray"))
}

func (suite *JSONWrapperTestSuite) TestCreateJSONArrayAtPath() {
	json, _ := New("{}")
	err := json.CreateJSONArrayAtPath("testArray")
	suite.Assert.Nil(err)
	suite.Assert.True(json.HasPath("testArray"))
}

func (suite *JSONWrapperTestSuite) TestArrayAppendInPath() {
	json, _ := New(`{"testArray": ["a", "b", "c"], "testString": "d"}`)
	err := json.ArrayAppendInPath("testArray",
		json.GetJSONObjectFromPath("testString"))
	suite.Assert.Nil(err)
}

func TestJSONWrapper(t *testing.T) {
	suite.Run(t, new(JSONWrapperTestSuite))
}
