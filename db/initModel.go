package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"linkshortener/log"
	"linkshortener/setting"
)

func InitModel() {
	NewModel(setting.Cfg.MongoDB.Database, "links")
	NewModel(setting.Cfg.MongoDB.Database, "link_access")

	statsIndex := mongo.IndexModel{
		Keys: bson.M{
			"hash": 1,
		},
		Options: options.Index().SetName("hash_index"),
	}

	statsTable := SetModel(setting.Cfg.MongoDB.Database, "link_access")
	err := statsTable.CreateOneIndex(statsIndex)
	if err != nil {
		log.PanicPrint("Failed to initialize database")
	}

}
