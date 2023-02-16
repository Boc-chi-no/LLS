package db

type MongoConnect struct {
	Name            string       `json:"name"`
	UserName        string       `json:"userName"`
	Password        string       `json:"password"`
	Hosts           []MongoHost  `json:"hosts"`
	Database        string       `json:"database"`
	Option          MongoOptions `json:"option"`
	ConnectTimeout  int          `json:"connectTimeout"`
	ExecuteTimeout  int          `json:"executeTimeout"`
	MinPoolSize     int          `json:"minPoolSize"`
	MaxPoolSize     int          `json:"maxPoolSize"`
	MaxConnIdleTime int          `json:"maxConnIdleTime"`
}

type MongoHost struct {
	Hst  string `json:"host"`
	Port string `json:"port"`
}

type MongoOptions struct {
	ReplicaSet       string `json:"replicaSet"`
	SlaveOk          bool   `json:"slaveOk"`
	Safe             bool   `json:"safe"`
	WtimeoutMS       int64  `json:"wtimeoutMS"`
	ConnectTimeoutMS int64  `json:"connectTimeoutMS"`
	SocketTimeoutMS  int64  `json:"socketTimeoutMS"`
	MaxPoolSize      int    `json:"maxPoolSize"`
	MinPoolSize      int    `json:"minPoolSize"`
	MaxIdleTimeMS    int64  `json:"maxIdleTimeMS"`
}
