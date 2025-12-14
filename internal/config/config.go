package config

type Config struct {
	HTTPAddr string
	GRPCAddr string
}

var cfg = &Config{
	HTTPAddr: ":8080",
	GRPCAddr: ":9090",
}

func LoadConfig() (*Config, error) {
	return cfg, nil
}
