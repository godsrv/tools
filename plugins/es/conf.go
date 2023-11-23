package es

type UserPassAuth struct {
	Username string
	Password string
}

type EsConf struct {
	Endpoints         []string
	Sniff             bool
	Gzip              bool
	EnableTrace       bool
	EnableHealthCheck bool
	Auth              bool
	UserPass          UserPassAuth
	Https             bool
}
