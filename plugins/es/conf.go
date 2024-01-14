package es

type UserPassAuth struct {
	Username string
	Password string
}

type EsConf struct {
	Endpoints         []string     `json:"endpoints"`
	Sniff             bool         `json:"sniff"`
	Gzip              bool         `json:"gzip"`
	EnableTrace       bool         `json:"enableTrace"`
	EnableHealthCheck bool         `json:"enableHealthCheck"`
	Auth              bool         `json:"auth"`
	UserPass          UserPassAuth `json:"user_pass"`
	Https             bool         `json:"https"`
}
