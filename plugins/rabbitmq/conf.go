package rabbitmq

type (
	// A RedisConf is a redis config.
	MQConf struct {
		Username string `json:",optional"`
		Password string `json:",optional"`
		Host     string `json:",optional"`
		Port     string `json:",optional"`
		Vhost    string `json:",optional"`
	}

	QueueDeclare struct {
		Name       string `json:",optional"`
		Durable    bool   `json:",optional"`
		AutoDelete bool   `json:",optional"`
		Exchange   string `json:",optional"`
	}
)
