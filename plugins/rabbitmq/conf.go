package rabbitmq

type (
	// A RedisConf is a redis config.
	MQConf struct {
		Username string `json:"username,optional"`
		Password string `json:"password,optional"`
		Host     string `json:"host,optional"`
		Port     string `json:"port,optional"`
		Vhost    string `json:"vhost,optional"`
	}

	QueueDeclare struct {
		Name       string `json:"name,optional"`
		Durable    bool   `json:"durable,optional"`
		AutoDelete bool   `json:"auto_delete,optional"`
		Exchange   string `json:"exchange,optional"`
	}
)
