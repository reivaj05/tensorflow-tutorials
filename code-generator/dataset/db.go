package GoDB

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var dbClient *DBClient

type DBClient struct {
	client *gorm.DB
}

type DBOptions struct {
	DBEngine   string
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
}

func Init(options *DBOptions, models ...interface{}) error {
	if !isDBEngineOK(options.DBEngine) {
		return fmt.Errorf("Wrong db engine")
	}
	if !isDBNameOK(options.DBName) {
		return fmt.Errorf("Wrong db name")
	}
	return newDBClient(options, models...)
}

func isDBEngineOK(engine string) bool {
	return engine != ""
}

func isDBNameOK(name string) bool {
	return name != ""
}

func GetDBClient() *DBClient {
	return dbClient
}

func newDBClient(options *DBOptions, models ...interface{}) error {
	client, err := gorm.Open(options.DBEngine, options.DBName)
	if err != nil {
		return err
	}
	dbClient = &DBClient{client: client}
	return dbClient.makeMigrations(models...)
}

func (db *DBClient) makeMigrations(models ...interface{}) error {
	return db.client.AutoMigrate(models...).Error
}

func (db *DBClient) Create(instanceModel interface{}) error {
	return db.client.Create(instanceModel).Error
}

func (db *DBClient) List(listModel interface{}) error {
	return db.client.Find(listModel).Error
}

func (db *DBClient) Get(instanceModel interface{}, id int) error {
	return db.client.First(instanceModel, id).Error
}

func (db *DBClient) GetWhere(instanceModel interface{}, key, value string) error {
	return db.client.Where(key+" = ?", value).First(instanceModel).Error
}

func (db *DBClient) Delete(instanceModel interface{}) error {
	return db.client.Delete(instanceModel).Error
}

func (db *DBClient) Update(instanceModel interface{}) error {
	return db.client.Save(instanceModel).Error
}
