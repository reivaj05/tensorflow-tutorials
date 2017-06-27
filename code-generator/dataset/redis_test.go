package GoRedis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	redisObj *RedisClient
}

func (suite *RedisTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
}

func (suite *RedisTestSuite) TearDownTest() {
	if suite.redisObj != nil {
		suite.redisObj.FlushDB()
	}
}

func (suite *RedisTestSuite) TestNewWithWrongAddress() {
	options := createRedisOptions("", "", 0)

	suite.assertNewWrongCase(options)
}

func (suite *RedisTestSuite) TestNewWithWrongPort() {
	options := createRedisOptions("localhost", "port", 0)

	suite.assertNewWrongCase(options)
}

func (suite *RedisTestSuite) TestNewWithWrongDBNumber() {
	options := createRedisOptions("localhost", "6379", -1)

	suite.assertNewWrongCase(options)
}

func (suite *RedisTestSuite) assertNewWrongCase(options *RedisOptions) {
	var err error
	suite.redisObj, err = New(options)
	suite.assert.Nil(suite.redisObj)
	suite.assert.NotNil(err)
}

func (suite *RedisTestSuite) TestInitWithGoodOptions() {
	var err error
	options := createRedisOptions("localhost", "6379", 5)
	suite.redisObj, err = New(options)
	suite.assert.NotNil(suite.redisObj)
	suite.assert.Nil(err)
}

func (suite *RedisTestSuite) TestSetValue() {
	options := createRedisOptions("localhost", "6379", 5)
	suite.redisObj, _ = New(options)
	key, value := "testKey", "testValue"
	err := suite.redisObj.SetValue(key, value, 0)
	suite.assert.Nil(err)
	result, err := suite.redisObj.GetValue(key)
	suite.assert.Nil(err)
	suite.assert.Equal(value, result)
}

func (suite *RedisTestSuite) TestGetValueNotFound() {
	options := createRedisOptions("localhost", "6379", 5)
	suite.redisObj, _ = New(options)
	key := "testKey"
	_, err := suite.redisObj.GetValue(key)
	suite.assert.NotNil(err)
}

func (suite *RedisTestSuite) TestExistsKeyTrue() {
	options := createRedisOptions("localhost", "6379", 5)
	suite.redisObj, _ = New(options)
	key, value := "testKey", "testValue"
	suite.redisObj.SetValue(key, value, 0)
	suite.assert.True(suite.redisObj.ExistsKey(key))
}

func (suite *RedisTestSuite) TestExistsKeyFalse() {
	options := createRedisOptions("localhost", "6379", 5)
	suite.redisObj, _ = New(options)
	key, _ := "testKey", "testValue"
	suite.assert.False(suite.redisObj.ExistsKey(key))
}

func createRedisOptions(address, port string, db int) *RedisOptions {
	return &RedisOptions{
		Address: address,
		Port:    port,
		DB:      db,
	}
}

// func (suite *RedisTestSuite) TestCSetValue() {
// 	assert := assert.New(suite.T())

// 	testKey := "testKeySet"
// 	testKey2 := "testKeySet2"
// 	testValue := "testValue"

// 	SetValue(testKey, testValue)
// 	assert.True(ExistsKey(testKey))
// 	assert.Equal(testValue, GetValue(testKey))

// 	Settings = &RedisSettings{DeleteAfterSeconds: 1}

// 	SetValue(testKey2, testValue)
// 	assert.True(ExistsKey(testKey2))
// 	time.Sleep(time.Second * 2)
// 	assert.False(ExistsKey(testKey2))
// }

// func (suite *RedisTestSuite) TestDGetValue() {
// 	assert := assert.New(suite.T())
// 	key, value, badKey := "testKeyGet", "testValue", "notExisting"
// 	SetValue(key, value)
// 	assert.Equal(value, GetValue(key))
// 	assert.Equal("", GetValue(badKey))
// }

// func (suite *RedisTestSuite) TestEExistsKey() {
// 	assert := assert.New(suite.T())
// 	key, value, badKey := "testKeyExists", "testValue", "notExisting"
// 	SetValue(key, value)
// 	assert.True(ExistsKey(key))
// 	assert.False(ExistsKey(badKey))
// }

// func (suite *RedisTestSuite) TestFAddMemberToSortedSet() {
// 	assert := assert.New(suite.T())
// 	key, badKey, score, member := "setKey", "badKey", 10.0, "member"
// 	count := AddMemberToSortedSet(score, key, member)
// 	assert.Equal(int64(1), count)

// 	SetValue(badKey, member)
// 	count = AddMemberToSortedSet(score, badKey, member)
// 	assert.Equal(int64(-1), count)
// }

// func (suite *RedisTestSuite) TestGDeleteMemberFromSortedSet() {
// 	assert := assert.New(suite.T())
// 	key, badKey, member, badMember := "setKey", "badKey", "member", "badMember"
// 	AddMemberToSortedSet(10.0, key, member)
// 	count := DeleteMemberFromSortedSet(key, member)
// 	assert.Equal(int64(1), count)

// 	count = DeleteMemberFromSortedSet(key, badMember)
// 	assert.Equal(int64(0), count)

// 	SetValue(badKey, member)
// 	count = DeleteMemberFromSortedSet(badKey, member)
// 	assert.Equal(int64(-1), count)
// }

// func (suite *RedisTestSuite) TestHExpireMembersBeforeNow() {
// 	assert := assert.New(suite.T())

// 	key, badKey, m1, m2, m3 := "setKey", "badKey", "member1", "member2", "member3"

// 	now := float64(time.Now().Add(time.Hour * 1).Unix())

// 	count := AddMemberToSortedSet(10, key, m1)
// 	count += AddMemberToSortedSet(100, key, m2)
// 	count += AddMemberToSortedSet(now, key, m3)

// 	res := len(GetMembersFromSortedSet(key))
// 	assert.Equal(int64(res), count)

// 	deleted := ExpireMembersBeforeNow(key)
// 	assert.Equal(int64(2), deleted)

// 	res = len(GetMembersFromSortedSet(key))
// 	assert.Equal(res, 1)
// 	deleted = ExpireMembersBeforeNow(key)
// 	assert.Equal(int64(0), deleted)

// 	SetValue(badKey, m1)
// 	deleted = ExpireMembersBeforeNow(badKey)
// 	assert.Equal(int64(-1), deleted)
// }

// func (suite *RedisTestSuite) TestIGetMembersFromSortedSet() {
// 	assert := assert.New(suite.T())

// 	key, m1, m2 := "setKey", "member1", "member2"

// 	res := len(GetMembersFromSortedSet(key))
// 	assert.Equal(0, res)

// 	count := AddMemberToSortedSet(10, key, m1)
// 	count += AddMemberToSortedSet(100, key, m2)

// 	res = len(GetMembersFromSortedSet(key))
// 	assert.Equal(int64(res), count)
// }

// func (suite *RedisTestSuite) TestJShouldIgnoreData() {
// 	assert := assert.New(suite.T())

// 	key, fields := "testKey", `{"TF1":"TF2","TF3":"TF4"}`

// 	ignore := ShouldIgnoreData(key, fields)
// 	assert.False(ignore)

// 	SetValue(key, fields)
// 	ignore = ShouldIgnoreData(key, fields)
// 	assert.True(ignore)

// 	fields = `{"TF5":"TF6"}`
// 	ignore = ShouldIgnoreData(key, fields)
// 	assert.False(ignore)
// }

// func (suite *RedisTestSuite) TestKGetFieldsForHash() {
// 	assert := assert.New(suite.T())

// 	json := `{
// 		"id": "id",
// 		"field1": "field1",
// 		"field2": "field2"
// 	}`

// 	doc, _ := rj.NewParsedStringJson(json)
// 	ct := doc.GetContainer()

// 	id, data := GetFields(ct)
// 	assert.Equal("id", id)
// 	assert.NotEqual("", data)
// }

// func (suite *RedisTestSuite) TestLHasChanged() {
// 	assert := assert.New(suite.T())

// 	key, json := "testKey", `{"value":"testHash"}`
// 	SetValue(key, json)
// 	assert.False(hasChanged(key, json))
// 	json = `{"value":"changed"}`
// 	assert.True(hasChanged(key, json))
// }

// func (suite *RedisTestSuite) TestNGetMembers() {
// 	assert := assert.New(suite.T())

// 	contains := func(members []string, member string) bool {
// 		for _, expected := range members {
// 			if member == expected {
// 				return true
// 			}
// 		}
// 		return false
// 	}

// 	key, badKey, members := "setKey", "badKey", []string{
// 		"member1",
// 		"member2",
// 		"member3",
// 	}

// 	size := len(GetMembers(key))
// 	assert.Equal(0, size)

// 	SetAdd(key, members)
// 	membersResult := GetMembers(key)
// 	assert.Equal(len(membersResult), len(members))
// 	for _, member := range membersResult {
// 		assert.True(contains(members, member))
// 	}

// 	SetValue(badKey, "test")
// 	membersResult = GetMembers(badKey)
// 	assert.Nil(membersResult)
// }

// func (suite *RedisTestSuite) TestOIssetMember() {
// 	assert := assert.New(suite.T())

// 	key, member := "setKey", "member"

// 	assert.False(IsSetMember(key, member))
// 	SetAdd(key, []string{member})
// 	assert.True(IsSetMember(key, member))
// }

// func (suite *RedisTestSuite) TestPSetAdd() {
// 	assert := assert.New(suite.T())

// 	key, badKey, members := "setKey", "badKey", []string{
// 		"member1",
// 		"member2",
// 		"member3",
// 	}
// 	count := SetAdd(key, members)
// 	assert.Equal(int64(len(members)), count)

// 	SetValue(badKey, "test")
// 	count = SetAdd(badKey, members)
// 	assert.Equal(int64(-1), count)
// }

// func (suite *RedisTestSuite) TestQSetRemove() {
// 	assert := assert.New(suite.T())

// 	key, badKey, members := "setKey", "badKey", []string{
// 		"member1",
// 		"member2",
// 		"member3",
// 	}
// 	count := SetRemove(key, members)
// 	assert.Equal(int64(0), count)

// 	SetAdd(key, members)
// 	count = SetRemove(key, members)
// 	assert.Equal(int64(len(members)), count)

// 	SetValue(badKey, "test")
// 	count = SetRemove(badKey, members)
// 	assert.Equal(int64(-1), count)
// }

func TestRedis(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
