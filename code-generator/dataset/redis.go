package GoRedis

import (
	"fmt"
	"strconv"
	"time"

	redis "gopkg.in/redis.v3"
)

type RedisClient struct {
	client *redis.Client
}

type RedisOptions struct {
	Address  string
	Port     string
	Password string
	DB       int
}

// New creates a redis client
func New(options *RedisOptions) (*RedisClient, error) {
	if !isAddressOK(options.Address) {
		return nil, fmt.Errorf("Wrong redis address")
	}
	if !isPortOK(options.Port) {
		return nil, fmt.Errorf("Wrong redis port")
	}
	if !isDBOK(options.DB) {
		return nil, fmt.Errorf("DB number can't be less than zero")
	}
	return createClient(options), nil
}

func isAddressOK(address string) bool {
	return address != ""
}

func isPortOK(port string) bool {
	_, err := strconv.Atoi(port)
	return err == nil
}

func isDBOK(db int) bool {
	return db >= 0
}

func createClient(options *RedisOptions) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     options.Address + ":" + options.Port,
			Password: options.Password,
			DB:       int64(options.DB),
		}),
	}
}

// SetValue sets a redis key with its value and expiration time
func (redisObj *RedisClient) SetValue(key, value string, expireTime int) error {
	t := time.Second * time.Duration(expireTime)
	_, err := redisObj.client.Set(key, value, t).Result()
	return err
}

// GetValue returns the value for the key passed as parameter
func (redisObj *RedisClient) GetValue(key string) (string, error) {
	return redisObj.client.Get(key).Result()
}

// ExistsKey returns true if the key exists in the db, false otherwise
func (redisObj *RedisClient) ExistsKey(key string) bool {
	result, _ := redisObj.client.Exists(key).Result()
	return result
}

func (redisObj *RedisClient) FlushDB() {
	redisObj.client.FlushDb()
}
