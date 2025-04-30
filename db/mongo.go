package db

import (
	"context"
	"errors"
	"fmt"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/setting"
	"time"

	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LlsMongoDB struct {
	DatabaseName string
	ConnectName  string
	dbPool       MongoPooler
}

// SetPool Setting up a connection pool
func (db *LlsMongoDB) SetPool(pool MongoPooler) *LlsMongoDB {
	db.dbPool = pool
	return db
}

// SetDB Set the connection name and database name
func (db *LlsMongoDB) SetDB(connectName, databaseName string) *LlsMongoDB {
	db.ConnectName = connectName
	db.DatabaseName = databaseName
	return db
}

// CreateConnectFunc Create a connection
type CreateConnectFunc func(*LlsMongoDB) []MongoConnect

func NewMongoDB() *LlsMongoDB {
	log.InfoPrint("Using the MongoDB as a data source")
	var mongoHosts []MongoHost
	if setting.Cfg.MongoDB.Cluster {
		for _, v := range setting.Cfg.MongoDB.IPs {
			mongoHosts = append(mongoHosts, MongoHost{
				Hst:  v,
				Port: setting.Cfg.MongoDB.Port,
			})
		}
	} else {
		mongoHosts = []MongoHost{
			{
				Hst:  setting.Cfg.MongoDB.IP,
				Port: setting.Cfg.MongoDB.Port,
			},
		}
	}

	var configs []MongoConnect
	configs = append(configs, MongoConnect{
		Name:            setting.Cfg.DB.Database,
		Database:        setting.Cfg.DB.Database,
		UserName:        setting.Cfg.MongoDB.User,
		Password:        setting.Cfg.MongoDB.Password,
		Hosts:           mongoHosts,
		ConnectTimeout:  setting.Cfg.MongoDB.ConnectTimeout,
		ExecuteTimeout:  setting.Cfg.MongoDB.ExecuteTimeout,
		MinPoolSize:     setting.Cfg.MongoDB.MinPoolSize,
		MaxPoolSize:     setting.Cfg.MongoDB.MaxPoolSize,
		MaxConnIdleTime: setting.Cfg.MongoDB.MaxConnIdleTime,
	})
	var db LlsMongoDB
	return db.SetPool(NewPools(configs))
}

type MongoDBTable struct {
	tableName string
	db        *LlsMongoDB
}

// NewMongoDBTable Initialization table
func NewMongoDBTable(db *LlsMongoDB, tableName string) Tabler {
	var table = &MongoDBTable{}
	table.tableName = tableName
	table.SetDB(db)
	table.CreateCollection()
	return table
}

// SetMongoDBTable Setting table
func SetMongoDBTable(db *LlsMongoDB, tableName string) Tabler {
	var table = &MongoDBTable{}
	table.tableName = tableName
	table.SetDB(db)
	return table
}

// CreateCollection Create a collection
func (t *MongoDBTable) CreateCollection() {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()

	collectionNames, err := db.Database.ListCollectionNames(ctx, bson.M{"name": t.tableName})
	if err != nil {
		log.WarnPrint("List Collection failed: %s", err)
	}

	for _, name := range collectionNames {
		if name == t.tableName {
			return
		}
	}

	err = db.Database.CreateCollection(ctx, t.tableName)
	if err != nil {
		log.WarnPrint("Create Collection failed: %s", err)
	}
}

// SetDB Setting up the database
func (t *MongoDBTable) SetDB(db interface{}) {
	mongoDB, ok := db.(*LlsMongoDB)
	if ok {
		t.db = mongoDB
	}
}

// getDB
func (t *MongoDBTable) getDB() *DatabaseInfo {
	return t.db.dbPool.GetDB(t.db.ConnectName)
}

func (t *MongoDBTable) CountDocuments(filter interface{}, _ *FindOptions) (int64, error) {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	count, err := db.Database.Collection(t.tableName).CountDocuments(ctx, filter)
	if err != nil {
		log.ErrorPrint("mongo CountDocuments error %v", err)
	}

	return count, err
}

func (t *MongoDBTable) InsertOne(document interface{}, autoKey bool) (interface{}, error) {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()

	doc := make(map[string]interface{})
	documentBson, _ := tool.MarshalJsonByBson(document)
	_ = json.Unmarshal(documentBson, &doc)
	id, ok := doc["_id"]

	if autoKey {
		if ok && fmt.Sprint(id) != "" {
			return nil, log.Errorf("_id should not be provided when autoKey is true")
		}
		delete(doc, "_id")
		document = doc
	} else {
		if !ok || fmt.Sprint(id) == "" {
			return nil, log.Errorf("_id is required when autoKey is false")
		}
	}

	result, err := db.Database.Collection(t.tableName).InsertOne(ctx, document)
	if err != nil {
		log.ErrorPrint("mongo InsertOne error %v", err)
	}
	return result.InsertedID, err
}

func (t *MongoDBTable) UpdateOne(filter interface{}, update interface{}) error {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	_, err := db.Database.Collection(t.tableName).UpdateOne(ctx, filter, update)
	if err != nil {
		log.ErrorPrint("mongo UpdateOne error %v", err)
	}
	return err
}

func (t *MongoDBTable) UpdateByID(id string, update interface{}) error {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	_, err := db.Database.Collection(t.tableName).UpdateByID(ctx, id, update)
	if err != nil {
		log.ErrorPrint("mongo UpdateByID error %v", err)
	}
	return err
}

func (t *MongoDBTable) FindOne(filter interface{}, result interface{}) error {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	err := db.Database.Collection(t.tableName).FindOne(ctx, filter).Decode(result)
	if err != nil {
		log.ErrorPrint("mongo FindOne error %v", err)
	}
	return err
}

func (t *MongoDBTable) FindByID(id interface{}, result interface{}) error {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	err := db.Database.Collection(t.tableName).FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.DebugPrint("No document found for id: %v", id)
			return err //TODO: 统一错误
		} else {
			log.ErrorPrint("MongoDB FindByID error %v", err)
			return err
		}
	}
	return nil
}

func (t *MongoDBTable) Find(filter interface{}, result interface{}, opt *FindOptions) error {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	cur, err := db.Database.Collection(t.tableName).Find(ctx, filter, options.Find().SetSkip(opt.Skip).SetLimit(opt.Limit).SetMin(opt.Min).SetMax(opt.Max))
	if err != nil {
		log.ErrorPrint("mongo Find error %v", err)
		return err

	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		_ = cur.Close(ctx)
	}(cur, context.Background())
	err = cur.All(context.Background(), result)
	if err != nil {
		log.ErrorPrint("mongo Find cur error %v", err)
	}
	return err
}

func (t *MongoDBTable) CreateOneIndex(indexInterface interface{}, opts ...interface{}) error {
	var createIndexesOptionsSlice []*options.CreateIndexesOptions

	for _, opt := range opts {
		if findOpt, ok := opt.(*options.CreateIndexesOptions); ok {
			createIndexesOptionsSlice = append(createIndexesOptionsSlice, findOpt)
		}
	}
	if index, ok := indexInterface.(mongo.IndexModel); ok {
		db := t.getDB()
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
		defer func() {
			cancel()
		}()
		_, err := db.Database.Collection(t.tableName).Indexes().CreateOne(ctx, index, createIndexesOptionsSlice...)
		if err != nil {
			log.ErrorPrint("mongo CreateOneIndex error %v", err)
		}
		return err
	}
	return fmt.Errorf("failed to type assertion failed: CreateOneIndex")
}
