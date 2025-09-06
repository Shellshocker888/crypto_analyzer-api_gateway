package model

type Config struct {
	Port                string
	AuthServiceURL      string
	PortfolioServiceURL string
	AlertServiceURL     string
	RedisCfg            *RedisConfig
}

type RedisConfig struct {
	Addr      string
	Password  string
	SessionDB int
}
