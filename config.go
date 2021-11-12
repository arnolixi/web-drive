package arc

type Config struct {
	*Server `mapstructure:"server"`
}

type Server struct {
	Addr string `mapstructure:"addr"`
}
