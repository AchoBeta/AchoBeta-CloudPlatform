package config

type Docker struct {
	Hub Hub `yaml:"hub" json:"hub"`
}

type Hub struct {
	Host     string `yaml:"host" json:"host"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}
