package db

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"linkshortener/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// configLock
var configLock sync.Mutex

type DatabaseInfo struct {
	Client   *mongo.Client
	Database *mongo.Database
	Config   MongoConnect
}

func (d *DatabaseInfo) Close() {
	if d.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if err := d.Client.Disconnect(ctx); err != nil {
			log.DebugPrint("mongodb close %v", err)
		}
	}
}

type MongoPooler interface {
	GetDB(connectName string) *DatabaseInfo
	Init(connect []MongoConnect)
	AddConnect(connect MongoConnect)
	AddConnects(connects []MongoConnect)
}

type MongoPool struct {
	pool map[string]*poolDB
}

func (m *MongoPool) GetDB(connectName string) *DatabaseInfo {
	return m.pool[connectName].dbInfo
}

func (m *MongoPool) Init(connects []MongoConnect) {
	m.pool = make(map[string]*poolDB)
	for _, connect := range connects {
		m.AddConnect(connect)
	}
}

// AddConnect 添加连接
func (m *MongoPool) AddConnect(connect MongoConnect) {
	configLock.Lock()
	defer configLock.Unlock()
	if _, ok := m.pool[connect.Name]; ok {
		log.PanicPrint("Mongo data connection already exists")
	}
	connectStr := m.open(connect)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(connect.ConnectTimeout)*time.Second)
	defer cancel()
	opts := options.Client().ApplyURI(connectStr)
	opts.SetMaxPoolSize(uint64(connect.MaxPoolSize))
	opts.SetMinPoolSize(uint64(connect.MinPoolSize))
	opts.SetMaxConnIdleTime(time.Duration(connect.MaxConnIdleTime) * time.Minute)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.DebugPrint("mongodb GetDB error %v", err)
	}
	database := client.Database(connect.Database)
	dbinfo := &DatabaseInfo{
		Client:   client,
		Database: database,
		Config:   connect,
	}

	if connectStr != "" {
		m.pool[connect.Name] = &poolDB{
			config:  connect,
			connect: connectStr,
			dbInfo:  dbinfo,
		}
	}
}

// AddConnects 添加多个连接
func (m *MongoPool) AddConnects(connects []MongoConnect) {
	for _, c := range connects {
		m.AddConnect(c)
	}
}

// open
func (m *MongoPool) open(connect MongoConnect) string {
	var userPassword string
	var hosts, connectOptions []string
	if connect.UserName != "" && connect.Password != "" {
		userPassword = fmt.Sprintf("%s:%s@", connect.UserName, connect.Password)
	}
	for _, h := range connect.Hosts {
		if h.Port != "" {
			hosts = append(hosts, fmt.Sprintf("%s:%s", h.Hst, h.Port))
		} else {
			hosts = append(hosts, fmt.Sprintf("%s", h.Hst))
		}

	}
	if connect.Option.ConnectTimeoutMS > 0 {
		connectOptions = append(connectOptions, fmt.Sprintf("connectTimeoutMS=%d", connect.Option.ConnectTimeoutMS))
	}
	if connect.Option.MaxIdleTimeMS > 0 {
		connectOptions = append(connectOptions, fmt.Sprintf("maxIdleTimeMS=%d", connect.Option.MaxIdleTimeMS))
	}
	if connect.Option.MaxPoolSize > 0 {
		connectOptions = append(connectOptions, fmt.Sprintf("maxPoolSize=%d", connect.Option.MaxPoolSize))
	}
	if connect.Option.MinPoolSize > 0 {
		connectOptions = append(connectOptions, fmt.Sprintf("minPoolSize=%d", connect.Option.MinPoolSize))
	}
	if connect.Option.SocketTimeoutMS > 0 {
		connectOptions = append(connectOptions, fmt.Sprintf("socketTimeoutMS=%d", connect.Option.SocketTimeoutMS))
	}
	if connect.Option.WtimeoutMS > 0 {
		connectOptions = append(connectOptions, fmt.Sprintf("wtimeoutMS=%d", connect.Option.WtimeoutMS))
	}
	if connect.Option.ReplicaSet != "" {
		connectOptions = append(connectOptions, fmt.Sprintf("replicaSet=%s", connect.Option.ReplicaSet))
	}
	connectOptions = append(connectOptions, fmt.Sprintf("safe=%t", connect.Option.Safe))
	connectOptions = append(connectOptions, fmt.Sprintf("slaveOk=%t", connect.Option.SlaveOk))
	return fmt.Sprintf("mongodb://%s%s/%s?%s", userPassword, strings.Join(hosts, ","), connect.Database, strings.Join(connectOptions, ","))
}

// poolDB
type poolDB struct {
	config  MongoConnect
	connect string
	dbInfo  *DatabaseInfo
}

// 全局实现者
var pool MongoPooler
var poolLock sync.Mutex

// NewPools 初始化多数据库连接
func NewPools(configs []MongoConnect) MongoPooler {
	poolLock.Lock()
	defer poolLock.Unlock()
	if pool == nil {
		pool = &MongoPool{}
		pool.Init(configs)
	}
	return pool
}

// NewPool 初始化数据库连接
//func NewPool(config MongoConnect) MongoPooler {
//	poolLock.Lock()
//	defer poolLock.Unlock()
//	if pool == nil {
//		var configs []MongoConnect
//		configs = append(configs, config)
//		pool = &MongoPool{}
//		pool.Init(configs)
//	}
//	return pool
//}
