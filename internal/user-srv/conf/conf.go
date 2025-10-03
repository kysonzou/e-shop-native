package conf

type Server_GRPC struct {
	Network  string `mapstructure:"network"`
	Addr     string `mapstructure:"addr"`
	// Timeout  int64  `mapstructure:"timeout"`
	// CertFile string `mapstructure:"cert_file"`
	// KeyFile  string `mapstructure:"key_file"`
}

type Server_HTTP struct {
	Network      string   `mapstructure:"network"`
	Addr         string   `mapstructure:"addr"`
	// Timeout      int64    `mapstructure:"timeout"`
	// CertFile     string   `mapstructure:"cert_file"`
	// KeyFile      string   `mapstructure:"key_file"`
	// CORS         bool     `mapstructure:"cors"`
	// AllowedHosts []string `mapstructure:"allowed_hosts"`
}

type Server_Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Server_MySQL struct {
	Drive        string `mapstructure:"drive"`
	DSN          string `mapstructure:"dsn"`
	MaxidleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxLifetime  int64  `mapstructure:"max_lifetime"`
}

type Auth struct {
	JwtKey string `mapstructure:"jwt_key"`
	Expire int64  `mapstructure:"expire"`
	Algorithm string `mapstructure:"algorithm"` 

}

type Log struct {
	Level string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}


type Server struct {
	HTTP *Server_HTTP `mapstructure:"http"`
	GRPC *Server_GRPC `mapstructure:"grpc"`
}

type Data struct {
	MySQL *Server_MySQL `mapstructure:"mysql"`
	Redis *Server_Redis `mapstructure:"redis"`
}

type Bootstrap struct {
	Server *Server `mapstructure:"server"`
	Data   *Data   `mapstructure:"data"`
	Auth   *Auth   `mapstructure:"auth"`
	Log    *Log    `mapstructure:"log"`
}
