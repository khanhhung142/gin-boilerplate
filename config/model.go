package config

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Log      LogConfig      `yaml:"log"`
	Auth     AuthConfig     `yaml:"auth"`
	Env      string         `yaml:"env"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	ConnectString string `yaml:"connect_string"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	Path       string `yaml:"path"`
	File       string `yaml:"file"`
	TimeFormat string `yaml:"timeFormat"`
}

type AuthConfig struct {
	SecretKey             string `yaml:"secret_key"`
	AccessTokenExpireTime int    `yaml:"access_token_expire_minute"`
}
