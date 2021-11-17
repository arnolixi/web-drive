package conf

import (
	"time"

	arc "gitee.com/Arno-lixi/web-drive"
)

type Config struct {
	*ServerConf `mapstructure:"server"`
	*DBConf     `mapstructure:"db"`
	*InfluxConf `mapstructure:"influxdb"`
}

func NewConfig() *Config {
	return &Config{
		NewServerConf(),
		NewDBConf(),
		NewInfluxConf(),
	}
}

func (c *Config) GetServer() arc.IServer {
	return c.ServerConf
}

func (c *Config) GetDB() arc.IDBConfig {
	return c.DBConf
}

func (c *Config) GetInfluxDB() arc.Influx {
	return c.InfluxConf
}

type ServerConf struct {
	ServerAddr     string        `mapstructure:"server_addr"`
	AuthPasswd     string        `mapstructure:"auth_passwd"`
	RegisterPasswd string        `mapstructure:"register_passwd"`
	HeartBeatTIme  time.Duration `mapstructure:"heart_beat_time"`
	Version        string        `mapstructure:"version"`
	Mode           string        `mapstructure:"mode"`
	NodeID         int64         `mapstructure:"node_id"`
	SignKey        string        `mapstructure:"sign_key"`
	RBACFile       string        `mapstructure:"rbac_file"`
}

func NewServerConf() *ServerConf {
	return &ServerConf{}
}

func (s *ServerConf) GetAddr() string {
	return s.ServerAddr
}

type DBConf struct {
	DriverName   string `mapstructure:"driver_name" json:"driver_name"`
	Dsn          string `mapstructure:"dsn" json:"dsn"`
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	Mode         string `mapstructure:"mode" json:"mode"`
}

func NewDBConf() *DBConf {
	return &DBConf{}
}

func (d *DBConf) GetDriverName() string {
	return d.DriverName
}

// 0 MaxIdleConns 1 MaxOpenConns
func (d *DBConf) GetConnSetting() []int {
	return []int{d.MaxIdleConns, d.MaxOpenConns}
}

func (d *DBConf) GetMode() string {
	return d.Mode
}

func (d *DBConf) GetDSN() string {
	return d.Dsn
}

type InfluxConf struct {
	Url      string `mapstructure:"url" json:"url"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Database string `mapstructure:"database" json:"database"`
}

func NewInfluxConf() *InfluxConf {
	return &InfluxConf{}
}

func (i *InfluxConf) GetUrl() string {
	return i.Url
}

func (i *InfluxConf) GetAuth() (username, passwd string) {
	return i.Username, i.Password
}

func (i *InfluxConf) GetDatabase() string {
	return i.Database
}
