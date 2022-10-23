package config

const (
	ADMIN_POW   = "MTQ1NacdNA"
	USER_POW    = "MjQ1NcctNA"
	TOURIST_POW = "MzQ1NzxzNA"
)

type Server struct {
	App    App    `yaml:"app" json:"app"`
	Db     Db     `yaml:"db" json:"db"`
	Docker Docker `yaml:"docker" json:"docker"`
}
