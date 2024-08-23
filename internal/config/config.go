package config

import (
	"github.com/spf13/viper"
	"strings"
)

const (
	defaultConfigFormat = "yaml"
	defaultConfigPath   = "./"
	defaultName         = "config.yaml"
	defaultEnvPrefix    = "EXPORTER"
)

type Config struct {
	path   string
	format string
	prefix string
	name   string

	LogLevel     string `mapstructure:"logLevel"`
	Registrators struct {
		Domains []string `mapstructure:"domains"`

		Godaddy struct {
			Secret string `mapstructure:"secret"`
			Token  string `mapstructure:"token"`
		} `mapstructure:"godaddy"`

		Pananame struct {
			Token string `mapstructure:"token"`
		} `mapstructure:"pananame"`
	} `mapstructure:"registrators"`
}

type Option func(*Config)

func WithName(name string) Option {
	return func(c *Config) {
		c.name = name
	}
}

func WithConfigFormat(format string) Option {
	return func(c *Config) {
		c.format = format
	}
}

func WithConfigPath(path string) Option {
	return func(c *Config) {
		c.path = path
	}
}

func New(options ...Option) (*Config, error) {
	cfg := &Config{
		path:   defaultConfigPath,
		format: defaultConfigFormat,
		prefix: defaultEnvPrefix,
		name:   defaultName,
	}

	// Apply options
	for _, opt := range options {
		opt(cfg)
	}

	viper.SetConfigName(cfg.name)
	viper.SetConfigType(cfg.format)
	viper.AddConfigPath(cfg.path)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(cfg.prefix)
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal config into struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
