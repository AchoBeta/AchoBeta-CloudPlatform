package config

const (
	ADMIN_POW   = "MTQ1NacdNA"
	USER_POW    = "MjQ1NcctNA"
	TOURIST_POW = "MzQ1NzxzNA"
)

type Server struct {
	App     App     `yaml:"app" json:"app"`
	Options Options `yaml:"options" json:"options"`
	Db      Db      `yaml:"db" json:"db"`
	Docker  Docker  `yaml:"docker" json:"docker"`
}

type App struct {
	Name      string `yaml:"name" json:"name"`
	Type      string `yaml:"type" json:"type"`
	Host      string `yaml:"host" json:"host"`
	Port      int    `yaml:"port" json:"port"`
	StartPort int    `yaml:"startPort" json:"startPort"`
	Salt      string `yaml:"salt" json:"salt"`
}

type Options struct {
	LogFilePath string `yaml:"logFilePath" json:"logFilePath"`
}

type Db struct {
	Mongo Mongo `yaml:"mongo" json:"mongo"`
	Redis Redis `yaml:"redis" json:"redis"`
}

type Mongo struct {
	Address    string `yaml:"address" json:"address"`
	Port       int    `yaml:"port" json:"port"`
	Username   string `yaml:"username" json:"username"`
	Password   string `yaml:"password" json:"password"`
	AuthSource string `yanl:"authSource" json:"authSource"`
}

type Redis struct {
	Address  string `yaml:"address" json:"address"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	Db       int    `yaml:"db" json:"db"`
}
