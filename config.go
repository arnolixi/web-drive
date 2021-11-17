package arc

type IConfig interface {
	GetServer() IServer
	GetDB() IDBConfig
	GetInfluxDB() Influx
	GetMode() string
}

type IDBConfig interface {
	GetDriverName() string
	GetDSN() string
	GetConnSetting() []int
}

type IServer interface {
	GetAddr() string
}

type Influx interface {
	GetUrl() string
	GetAuth() (username, passwd string)
	GetDatabase() string
}
