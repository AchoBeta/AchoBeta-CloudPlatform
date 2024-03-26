package config

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
