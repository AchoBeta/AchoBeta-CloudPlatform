package config

type App struct {
	Name      string `yaml:"name" json:"name"`
	Engine    string `yaml:"engine" json:"engine"`
	Host      string `yaml:"host" json:"host"`
	Port      int    `yaml:"port" json:"port"`
	StartPort int    `yaml:"startPort" json:"startPort"`
	Salt      string `yaml:"salt" json:"salt"`
	Lark      Lark   `yaml:"lark" json:"lark"`
}
