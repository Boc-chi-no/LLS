package model

import "os"

type Config struct {
	RunMode   string `ini:"RUN_MODE"`

	LOG             LOGConfig          `ini:"log"`
	HTTP            HTTPConfig         `ini:"http"`
	HTTPLimiter     HTTPLimiterConfig  `ini:"http_limiter"`
	MongoDB         MongoDBConfig      `ini:"mongodb"`
}

type LOGConfig struct {
	Debug               bool    `ini:"DEBUG"`
	File *os.File
}

type HTTPConfig struct {
	Listen               string `ini:"LISTEN"`
	RandomSessionSecret  bool   `ini:"RANDOM_SESSION_SECRET"`
	SessionSecret        string `ini:"SESSION_SECRET"`
	FilesEmbed           bool   `ini:"STATIC_FILES_EMBED"`
	FilesURI             string `ini:"STATIC_FILES_URI"`
}

type HTTPLimiterConfig struct {
	EnableLimiter     bool   `ini:"ENABLE_LIMITER"`
	LimitRate         int    `ini:"LIMIT_RATE"`
	LimitBurst        int    `ini:"LIMIT_BURST"`
	Timeout           int    `ini:"TIMEOUT"`
}

type MongoDBConfig struct {
	IP               string    `ini:"IP"`
	Cluster          bool      `ini:"CLUSTER"`
	IPs              []string  `ini:"IPS"`
	Port             string    `ini:"PORT"`
	User             string    `ini:"USER"`
	Password         string    `ini:"PASSWORD"`
	Database         string    `ini:"DATABASE"`
	ConnectTimeout   int       `ini:"CONNECT_TIMEOUT"`
	ExecuteTimeout   int       `ini:"EXECUTE_TIMEOUT"`
	MinPoolSize      int       `ini:"MIN_POOL_SIZE"`
	MaxPoolSize      int       `ini:"MAX_POOL_SIZE"`
	MaxConnIdleTime  int       `ini:"MAX_CONN_IDLE_TIME"`
}
