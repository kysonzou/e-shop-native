package conf

type GRPC struct {
	Network  string  `json:"network"`
	Address  string  `json:"address"`
	Timeout  int64   `json:"timeout"`
	CertFile string  `json:"cert_file"`
	KeyFile  string  `json:"key_file"`
}

type HTTP struct {
	Network      string `json:"network"`
	Address      string `json:"address"`
	Timeout      int64  `json:"timeout"`
	CertFile     string `json:"cert_file"`
	KeyFile      string `json:"key_file"`
	CORS         bool   `json:"cors"`
	AllowedHosts []string `json:"allowed_hosts"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type MySQL struct {
	Drive string `json:"drive"`
	DSN   string	`json:"dsn"`
	MaxidleConns int `json:"max_idle_conns"`
	MaxOpenConns int `json:"max_open_conns"`
	MaxLifetime  int64 `json:"max_lifetime"`
}

type Server struct{
	HTTP *HTTP `json:"http"`
	GRPC *GRPC `json:"grpc"`
}

type Data struct {
	MySQL *MySQL `json:"mysql"`
	Redis *Redis `json:"redis"`
}


type Bootstrap struct {
	Server *Server  `json:"server"`
	Data *Data `json:"data"`
}
