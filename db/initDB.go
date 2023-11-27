package db

import (
	"linkshortener/log"
	"linkshortener/setting"
	"strings"
)

var MongoDB *LlsMongoDB
var BadgerDB *LlsBadgerDB

type Tabler interface {
	SetDB(db interface{})
	InsertOne(document interface{}, key string, autoKey bool) error
	UpdateOne(filter interface{}, result interface{}) error
	UpdateByID(id string, update interface{}) error
	FindByID(id interface{}, result interface{}) error
	FindOne(filter interface{}, result interface{}) error
	Find(filter interface{}, result interface{}, opts *FindOptions) error
	CreateOneIndex(index interface{}, opts ...interface{}) error
	CountDocuments(filter interface{}, opt *FindOptions) (int64, error)
}

func NewModel(dbName, tableName string) Tabler {
	switch strings.ToUpper(setting.Cfg.DB.Type) {
	case "BADGERDB":
		return NewBadgerDBTable(BadgerDB.SetDB(dbName, dbName), tableName)
	case "MONGODB":
		return NewMongoDBTable(MongoDB.SetDB(dbName, dbName), tableName)
	default:
		return nil
	}

}

func SetModel(dbName, tableName string) Tabler {
	switch strings.ToUpper(setting.Cfg.DB.Type) {
	case "BADGERDB":
		return SetBadgerDBTable(BadgerDB.SetDB(dbName, dbName), tableName)
	case "MONGODB":
		return SetMongoDBTable(MongoDB.SetDB(dbName, dbName), tableName)
	default:
		return nil
	}
}

func InitDB() {
	switch strings.ToUpper(setting.Cfg.DB.Type) {
	case "BADGERDB":
		BadgerDB = NewBadgerDB()
	case "MONGODB":
		MongoDB = NewMongoDB()
	default:
		log.PanicPrint("Database types are only allowed to be BadgerDB|MongoDB")
	}
}
