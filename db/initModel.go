package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"linkshortener/setting"
)

func InitModel()  {
	statsIndex := mongo.IndexModel{
		Keys: bson.M{
			"hash": 1,
		},
		Options: options.Index().SetName("hash_index"),
	}
	statsTable:= NewModel(setting.Cfg.MongoDB.Database, "link_access")
	statsTable.CreateOneIndex(statsIndex)
}