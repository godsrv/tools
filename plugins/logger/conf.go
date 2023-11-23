package logger

type ZapConf struct {
	Level     string `json:"level"`
	Director  string `json:"director"`
	MaxSize   int    `json:"maxSize"`
	MaxAge    int    `json:"maxAge"`
	MaxBackup int    `json:"maxBackup"`
}
