package config

type Lark struct {
	AppId       string `yaml:"appId" json:"appId"`
	AppSecret   string `yaml:"appSecret" json:"appSecret"`
	RedirectUrl string `yaml:"redirectUrl" json:"redirectUrl"`
}
