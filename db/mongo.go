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

//MongoDB
type MongoDB struct {
	DatabaseName string
	ConnectName  string
	dbPool       MongoPooler
}

//SetPool 设置连接池
func (db *MongoDB)SetPool(pool MongoPooler)*MongoDB{
	db.dbPool=pool
	return db
}

//SetDB 设置连接名和数据库名称
func (db *MongoDB)SetDB(connectName,databaseName string)*MongoDB{
	db.ConnectName =connectName
	db.DatabaseName =databaseName
	return db
}

//CreateConnectFunc 创建连接
type CreateConnectFunc func(*MongoDB)[]MongoConnect

//InjectCreatePool 注入连接池的创建过程
//connect 连接mongo配置 可以为空,如果为空需要自己创建数据连接查询配置后返回
//fun 回调函数返回mongo连接配置
func InjectCreatePool(fun CreateConnectFunc){
	if fun!=nil {
		configs:= fun(db)
		if len(configs)>0 {
			db.dbPool.AddConnects(configs)
		}
	}
}

//NewDB 
func NewDB()*MongoDB{

	var mongoHosts  []MongoHost
	if setting.Cfg.MongoDB.Cluster{
		for _, v := range setting.Cfg.MongoDB.IPs {
			mongoHosts = append(mongoHosts, MongoHost{
				Hst:  v,
				Port: setting.Cfg.MongoDB.Port,
			})
		}
	}else {
		mongoHosts = []MongoHost{
			{
				Hst:  setting.Cfg.MongoDB.IP,
				Port: setting.Cfg.MongoDB.Port,
			},
		}
	}

	var configs []MongoConnect
	configs=append(configs,MongoConnect{
		Name: setting.Cfg.MongoDB.Database,
		Database: setting.Cfg.MongoDB.Database,
		UserName: setting.Cfg.MongoDB.User,
		Password: setting.Cfg.MongoDB.Password,
		Hosts: mongoHosts,
		ConnectTimeout: setting.Cfg.MongoDB.ConnectTimeout,
		ExecuteTimeout: setting.Cfg.MongoDB.ExecuteTimeout,
		MinPoolSize: setting.Cfg.MongoDB.MinPoolSize,
		MaxPoolSize:setting.Cfg.MongoDB.MaxPoolSize,
		MaxConnIdleTime: setting.Cfg.MongoDB.MaxConnIdleTime,
	})
	var db MongoDB
	return db.SetPool(NewPools(configs))
}

type Tabler interface{
	SetDB(db *MongoDB)
	InsertOne(document interface{}) *mongo.InsertOneResult
	UpdateOne(filter interface{},result interface{}) *mongo.UpdateResult
	UpdateByID(id interface{},result interface{}) *mongo.UpdateResult
	FindByID(id interface{},result interface{})
	FindOne(filter interface{},result interface{})
	Find(filter interface{},result interface{},opts ...*options.FindOptions)
	CreateOneIndex(index mongo.IndexModel, opts ...*options.CreateIndexesOptions)
	CountDocuments(filter interface{}) int64
}

//Table 操作表
type Table struct{
	tableName string
	db *MongoDB
}

//NewTable 初始化表
func NewTable(db *MongoDB,tableName string) Tabler{
	var table *Table=&Table{}
	table.tableName=tableName
	table.SetDB(db)
	table.CreateCollection()
	return table
}
//CreateCollection 创建集合
func (t *Table)CreateCollection(){
	db:=t.getDB()
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func(){
		cancel()
	}()
	db.Database.CreateCollection(ctx, t.tableName)
}

//SetDB 设置数据库
func (t *Table) SetDB(db *MongoDB){
	t.db=db
}

//getDB
func (t *Table)getDB() *DBInfo{
	return t.db.dbPool.GetDB(t.db.ConnectName)
}

func (t *Table)CountDocuments(filter interface{}) int64{
	db:=t.getDB()
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func(){
		cancel()
	}()
	count, err := db.Database.Collection(t.tableName).CountDocuments(ctx, filter)
	if err!=nil{
		log.ErrorPrint("mongo CountDocuments error %v", err)
	}

	return count
}

func (t *Table)InsertOne(document interface{})*mongo.InsertOneResult{
	db:=t.getDB()
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func(){
		cancel()
	}()
	res,err:=db.Database.Collection(t.tableName).InsertOne(ctx,document)
	if err!=nil{
		log.ErrorPrint("mongo InsertOne error %v", err)
	}

	return res
}

func (t *Table)UpdateOne(filter interface{}, update interface{}) *mongo.UpdateResult {
	db:=t.getDB()
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func(){
		cancel()
	}()
	res,err:=db.Database.Collection(t.tableName).UpdateOne(ctx,filter,update)
	if err!=nil{
		log.ErrorPrint("mongo UpdateOne error %v", err)
	}
	return res
}

func (t *Table)UpdateByID(id interface{}, update interface{}) *mongo.UpdateResult {
	db:=t.getDB()
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func(){
		cancel()
	}()
	res,err:=db.Database.Collection(t.tableName).UpdateByID(ctx,id,update)
	if err!=nil{
		log.ErrorPrint("mongo UpdateByID error %v", err)
	}
	return res
}

func (t *Table)FindOne(filter interface{},result interface{}){
	db:=t.getDB()
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func(){
		cancel()
	}()
	err:=db.Database.Collection(t.tableName).FindOne(ctx,filter).Decode(result)
	if err!=nil{
		log.ErrorPrint("mongo FindOne error %v", err)
	}
}

func (t *Table)FindByID(id interface{},result interface{}){
	db:=t.getDB()
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func(){
		cancel()
	}()
	cur,err:=db.Database.Collection(t.tableName).Find(ctx,bson.D{{Key:"_id",Value: id}})

	if err!=nil{
		log.ErrorPrint("mongo Find error %v", err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		_ = cur.Close(ctx)
	}(cur, context.Background())
	err=cur.All(context.Background(), result)

	if err!=nil{
		log.ErrorPrint("mongo Find cur error %v", err)
	}
}


func (t *Table)Find(filter interface{},result interface{},opts ...*options.FindOptions){
	db:=t.getDB()
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func(){
		cancel()
	}()
	cur,err:=db.Database.Collection(t.tableName).Find(ctx,filter,opts...)
	if err!=nil{
		log.ErrorPrint("mongo Find error %v", err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		_ = cur.Close(ctx)
	}(cur, context.Background())
	err=cur.All(context.Background(), result)
	if err!=nil{
		log.ErrorPrint("mongo Find cur error %v", err)
	}
}

func (t *Table)CreateOneIndex(index mongo.IndexModel, opts ...*options.CreateIndexesOptions){
	db:=t.getDB()
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(db.Config.ExecuteTimeout)*time.Second)
	defer func(){
		cancel()
	}()
	_,err:=db.Database.Collection(t.tableName).Indexes().CreateOne(ctx,index,opts...)
	if err!=nil{
		log.ErrorPrint("mongo CreateOneIndex error %v", err)
	}
}


func NewModel(dbName,tableName string)Tabler{
	return NewTable(db.SetDB(dbName, dbName), tableName)
}

//db 全局变量
var db *MongoDB

func InitDB(){
	db=NewDB()

}