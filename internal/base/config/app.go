package config

type App struct {
	Name      string `yaml:"name" json:"name"`
	Type      string `yaml:"type" json:"type"`
	Host      string `yaml:"host" json:"host"`
	Port      int    `yaml:"port" json:"port"`
	StartPort int    `yaml:"startPort" json:"startPort"`
	Salt      string `yaml:"salt" json:"salt"`
}
