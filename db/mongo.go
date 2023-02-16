package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"linkshortener/log"
	"linkshortener/setting"
	"time"
)

type MongoDB struct {
	DatabaseName string
	ConnectName  string
	dbPool       MongoPooler
}

// SetPool 设置连接池
func (db *MongoDB) SetPool(pool MongoPooler) *MongoDB {
	db.dbPool = pool
	return db
}

// SetDB 设置连接名和数据库名称
func (db *MongoDB) SetDB(connectName, databaseName string) *MongoDB {
	db.ConnectName = connectName
	db.DatabaseName = databaseName
	return db
}

// CreateConnectFunc 创建连接
type CreateConnectFunc func(*MongoDB) []MongoConnect

func NewDB() *MongoDB {

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
		Name:            setting.Cfg.MongoDB.Database,
		Database:        setting.Cfg.MongoDB.Database,
		UserName:        setting.Cfg.MongoDB.User,
		Password:        setting.Cfg.MongoDB.Password,
		Hosts:           mongoHosts,
		ConnectTimeout:  setting.Cfg.MongoDB.ConnectTimeout,
		ExecuteTimeout:  setting.Cfg.MongoDB.ExecuteTimeout,
		MinPoolSize:     setting.Cfg.MongoDB.MinPoolSize,
		MaxPoolSize:     setting.Cfg.MongoDB.MaxPoolSize,
		MaxConnIdleTime: setting.Cfg.MongoDB.MaxConnIdleTime,
	})
	var db MongoDB
	return db.SetPool(NewPools(configs))
}

type Tabler interface {
	SetDB(db *MongoDB)
	InsertOne(document interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(filter interface{}, result interface{}) (*mongo.UpdateResult, error)
	UpdateByID(id interface{}, result interface{}) (*mongo.UpdateResult, error)
	FindByID(id interface{}, result interface{}) error
	FindOne(filter interface{}, result interface{}) error
	Find(filter interface{}, result interface{}, opts ...*options.FindOptions) error
	CreateOneIndex(index mongo.IndexModel, opts ...*options.CreateIndexesOptions) error
	CountDocuments(filter interface{}) (int64, error)
}

// Table 操作表
type Table struct {
	tableName string
	db        *MongoDB
}

// NewTable 初始化表
func NewTable(db *MongoDB, tableName string) Tabler {
	var table = &Table{}
	table.tableName = tableName
	table.SetDB(db)
	table.CreateCollection()
	return table
}

// SetTable 设置表
func SetTable(db *MongoDB, tableName string) Tabler {
	var table = &Table{}
	table.tableName = tableName
	table.SetDB(db)
	return table
}

// CreateCollection 创建集合
func (t *Table) CreateCollection() {
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

// SetDB 设置数据库
func (t *Table) SetDB(db *MongoDB) {
	t.db = db
}

// getDB
func (t *Table) getDB() *DatabaseInfo {
	return t.db.dbPool.GetDB(t.db.ConnectName)
}

func (t *Table) CountDocuments(filter interface{}) (int64, error) {
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

func (t *Table) InsertOne(document interface{}) (*mongo.InsertOneResult, error) {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	res, err := db.Database.Collection(t.tableName).InsertOne(ctx, document)
	if err != nil {
		log.ErrorPrint("mongo InsertOne error %v", err)
	}

	return res, err
}

func (t *Table) UpdateOne(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	res, err := db.Database.Collection(t.tableName).UpdateOne(ctx, filter, update)
	if err != nil {
		log.ErrorPrint("mongo UpdateOne error %v", err)
	}
	return res, err
}

func (t *Table) UpdateByID(id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	res, err := db.Database.Collection(t.tableName).UpdateByID(ctx, id, update)
	if err != nil {
		log.ErrorPrint("mongo UpdateByID error %v", err)
	}
	return res, err
}

func (t *Table) FindOne(filter interface{}, result interface{}) error {
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

func (t *Table) FindByID(id interface{}, result interface{}) error {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	cur, err := db.Database.Collection(t.tableName).Find(ctx, bson.D{{Key: "_id", Value: id}})

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

func (t *Table) Find(filter interface{}, result interface{}, opts ...*options.FindOptions) error {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	cur, err := db.Database.Collection(t.tableName).Find(ctx, filter, opts...)
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

func (t *Table) CreateOneIndex(index mongo.IndexModel, opts ...*options.CreateIndexesOptions) error {
	db := t.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	_, err := db.Database.Collection(t.tableName).Indexes().CreateOne(ctx, index, opts...)
	if err != nil {
		log.ErrorPrint("mongo CreateOneIndex error %v", err)
	}
	return err
}

func NewModel(dbName, tableName string) Tabler {
	return NewTable(db.SetDB(dbName, dbName), tableName)
}

func SetModel(dbName, tableName string) Tabler {
	return SetTable(db.SetDB(dbName, dbName), tableName)
}

// db 全局变量
var db *MongoDB

func InitDB() {
	db = NewDB()

}
