package model

import "os"

type Config struct {
	RunMode string `ini:"RUN_MODE"`
	Seed    uint32 `ini:"GENERATE_SEED"`

	LOG         LOGConfig         `ini:"log"`
	I18N        I18NConfig        `ini:"i18n"`
	HTTP        HTTPConfig        `ini:"http"`
	HTTPLimiter HTTPLimiterConfig `ini:"http_limiter"`
	MongoDB     MongoDBConfig     `ini:"mongodb"`
}

type LOGConfig struct {
	Debug bool `ini:"DEBUG"`
	File  *os.File
}

type I18NConfig struct {
	AddExtraLanguage   bool   `ini:"ADD_EXTRA_LANGUAGE"`
	ExtraLanguageName  string `ini:"EXTRA_LANGUAGE_NAME"`
	ExtraLanguageFiles string `ini:"EXTRA_LANGUAGE_FILES"`
}

type HTTPConfig struct {
	Listen               string `ini:"LISTEN"`
	BasePath             string `ini:"BASE_PATH"`
	SoftRedirectBasePath string `ini:"SOFT_REDIRECT_BASE_PATH"`
	RandomSessionSecret  bool   `ini:"RANDOM_SESSION_SECRET"`
	SessionSecret        string `ini:"SESSION_SECRET"`
	DisableFilesDirEmbed bool   `ini:"DISABLE_STATIC_FILES_DIR_EMBED"`
	FilesDirURI          string `ini:"STATIC_FILES_DIR_URI"`
}

type HTTPLimiterConfig struct {
	EnableLimiter bool `ini:"ENABLE_LIMITER"`
	LimitRate     int  `ini:"LIMIT_RATE"`
	LimitBurst    int  `ini:"LIMIT_BURST"`
	Timeout       int  `ini:"TIMEOUT"`
}

type MongoDBConfig struct {
	IP              string   `ini:"IP"`
	Cluster         bool     `ini:"CLUSTER"`
	IPs             []string `ini:"IPS"`
	Port            string   `ini:"PORT"`
	User            string   `ini:"USER"`
	Password        string   `ini:"PASSWORD"`
	Database        string   `ini:"DATABASE"`
	ConnectTimeout  int      `ini:"CONNECT_TIMEOUT"`
	ExecuteTimeout  int      `ini:"EXECUTE_TIMEOUT"`
	MinPoolSize     int      `ini:"MIN_POOL_SIZE"`
	MaxPoolSize     int      `ini:"MAX_POOL_SIZE"`
	MaxConnIdleTime int      `ini:"MAX_CONN_IDLE_TIME"`
}
