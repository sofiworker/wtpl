package conf

type Config struct {
	App    App    `yaml:"app"`
	Logger Logger `yaml:"logger"`
	Server Server `yaml:"server"`
	Rpc    Rpc    `yaml:"rpc"`
}

type App struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

type Logger struct {
	Path     []string `yaml:"path"`
	Level    string   `yaml:"level"`
	Compress bool     `yaml:"compress"`
	Retain   string   `yaml:"retain"`
	Filesize string   `yaml:"filesize"`
}

type Server struct {
	Ip      string `yaml:"ip"`
	Port    int    `yaml:"port"`
	TLSKey  string `yaml:"tlsKey"`
	TLSCert string `yaml:"tlsCert"`
}

type Rpc struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	TLSKey  string `yaml:"tlsKey"`
	TLSCert string `yaml:"tlsCert"`
}
