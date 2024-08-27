package conf

type Bootstrap struct {
	Server Server `json:"server,omitempty"`
	Log    Log    `json:"log,omitempty"`
}

type Server struct {
	Name    string     `json:"name,omitempty"`
	Version string     `json:"version,omitempty"`
	HTTP    HTTPServer `json:"http,omitempty"`
	GRPC    GRPCServer `json:"grpc,omitempty"`
}

type HTTPServer struct {
	Network string `json:"network,omitempty"`
	Addr    string `json:"addr,omitempty"`
	Timeout int32  `json:"timeout,omitempty"`
}

type GRPCServer struct {
	Network string `json:"network,omitempty"`
	Addr    string `json:"addr,omitempty"`
	Timeout int32  `json:"timeout,omitempty"`
}

type Log struct {
	MaxSize    int  `json:"max_size,omitempty"`
	MaxBackups int  `json:"max_backups,omitempty"`
	MaxAge     int  `json:"max_age,omitempty"`
	Compress   bool `json:"compress,omitempty"`
	LocalTime  bool `json:"local_time,omitempty"`
}
