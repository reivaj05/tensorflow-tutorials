package GoDB

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

type testModel struct {
	gorm.Model
	test int
}

func (suite *DBTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
}

func (suite *DBTestSuite) SetupTest() {
	dbClient = nil
}

func (suite *DBTestSuite) TearDownTest() {
	os.Remove("test.db")
}

func (suite *DBTestSuite) TestSuccessfulInit() {
	options := createOptions("sqlite3", "test.db")
	suite.assert.Nil(Init(options))
}

func (suite *DBTestSuite) TestInitWithoutDBEngine() {
	options := createOptions("", "test.db")
	suite.assert.NotNil(Init(options))
}

func (suite *DBTestSuite) TestInitWithoutDBName() {
	options := createOptions("sqlite3", "")
	suite.assert.NotNil(Init(options))
}

func (suite *DBTestSuite) TestGetDBClientWithoutInit() {
	client := GetDBClient()
	suite.assert.Nil(client)
}

func (suite *DBTestSuite) TestGetDBClientWithInit() {
	options := createOptions("sqlite3", "test.db")
	Init(options)
	client := GetDBClient()
	suite.assert.NotNil(client)
}

func (suite *DBTestSuite) TestCreateWithWrongModelInstance() {
	client := initAndGetClient(nil)
	err := client.Create(nil)
	suite.assert.NotNil(err)
}

func (suite *DBTestSuite) TestCreateWithGoodModelInstance() {
	client := initAndGetClient([]interface{}{&testModel{}}...)
	test := testModel{}
	err := client.Create(&test)
	suite.assert.Nil(err)
}

func (suite *DBTestSuite) TestListWithWrongModelInstance() {
	client := initAndGetClient(nil)
	var elements []testModel
	err := client.List(&elements)
	suite.assert.NotNil(err)
}

func (suite *DBTestSuite) TestListWithGoodModelInstance() {
	client := initAndGetClient([]interface{}{&testModel{}}...)
	var elements []testModel
	err := client.List(&elements)
	suite.assert.Nil(err)
}

func (suite *DBTestSuite) TestGetWithWrongModelInstance() {
	client := initAndGetClient([]interface{}{&testModel{}}...)
	test := testModel{}
	err := client.Get(&test, 1)
	suite.assert.NotNil(err)
}

func (suite *DBTestSuite) TestGetWithGoodModelInstance() {
	client := initAndGetClient([]interface{}{&testModel{}}...)
	test := testModel{}
	client.Create(&test)
	err := client.Get(&test, 1)
	suite.assert.Nil(err)
}

func (suite *DBTestSuite) TestGetWhereWithWrongModelInstance() {
	client := initAndGetClient([]interface{}{&testModel{}}...)
	test := testModel{}
	err := client.GetWhere(&test, "id", "1")
	suite.assert.NotNil(err)
}

func (suite *DBTestSuite) TestGetWhereWithGoodModelInstance() {
	client := initAndGetClient([]interface{}{&testModel{}}...)
	test := testModel{}
	client.Create(&test)
	err := client.GetWhere(&test, "id", "1")
	suite.assert.Nil(err)
}

func (suite *DBTestSuite) TestDeleteWithWrongModelInstance() {
	client := initAndGetClient(nil)
	test := testModel{}
	err := client.Delete(&test)
	suite.assert.NotNil(err)
}

func (suite *DBTestSuite) TestDeleteWithGoodModelInstance() {
	client := initAndGetClient([]interface{}{&testModel{}}...)
	test := testModel{}
	client.Create(&test)
	err := client.Delete(&test)
	suite.assert.Nil(err)
}

func (suite *DBTestSuite) TestUpdateWithWrongModelInstance() {
	client := initAndGetClient(nil)
	test := testModel{}
	err := client.Update(&test)
	suite.assert.NotNil(err)
}

func (suite *DBTestSuite) TestUpdateWithGoodModelInstance() {
	client := initAndGetClient([]interface{}{&testModel{}}...)
	test := testModel{}
	client.Create(&test)
	test.test = 10
	err := client.Update(&test)
	suite.assert.Nil(err)
}

func initAndGetClient(models ...interface{}) *DBClient {
	options := createOptions("sqlite3", "test.db")
	Init(options, models...)
	return GetDBClient()
}

func createOptions(engine, name string) *DBOptions {
	return &DBOptions{
		DBEngine: engine,
		DBName:   name,
	}
}

func TestDB(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
