package db

import (
	"linkshortener/log"
	"linkshortener/setting"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitModel() {
	//linksTable := NewModel(setting.Cfg.DB.Database, "links")
	statsTable := NewModel(setting.Cfg.DB.Database, "link_access")
	switch strings.ToUpper(setting.Cfg.DB.Type) {
	case "BADGERDB":
		statsIndex := mongo.IndexModel{
			Keys: bson.M{
				"hash": 1,
			},
			Options: options.Index().SetName("hash_index"),
		}
		err := statsTable.CreateOneIndex(statsIndex)
		if err != nil {
			log.PanicPrint("Failed to initialize MongoDB")
		}
	case "MONGODB":
		statsIndex := mongo.IndexModel{
			Keys: bson.M{
				"hash": 1,
			},
			Options: options.Index().SetName("hash_index"),
		}
		err := statsTable.CreateOneIndex(statsIndex)
		if err != nil {
			log.PanicPrint("Failed to initialize MongoDB")
		}
	default:
		return
	}

}
