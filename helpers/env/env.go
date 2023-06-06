package env

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresHost string `mapstructure:"POSTGRES_HOST" default:""`
	PostgresPort int    `mapstructure:"POSTGRES_PORT" default:""`
	PostgresDB   string `mapstructure:"POSTGRES_DB" default:""`
	PostgresUser string `mapstructure:"POSTGRES_USER" default:""`
	PostgresPass string `mapstructure:"POSTGRES_PASSWORD" default:""`
	AppHost      string `mapstructure:"APP_HOST" default:""`
	AppPort      int    `mapstructure:"APP_PORT" default:""`
	PostgresSSL  string `mapstructure:"POSTGRES_SSL" default:"disable"`
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=Asia/Jakarta",
		c.PostgresHost, c.PostgresUser, c.PostgresPass, c.PostgresDB, c.PostgresPort, c.PostgresSSL)
}

var (
	cfg *Config = nil
)

func Get() *Config {
	if cfg == nil {
		cfg = new(Config)

		viper.SetConfigType("yaml")
		viper.SetConfigFile("env.yaml")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		_ = viper.ReadInConfig()

		e := reflect.ValueOf(cfg).Elem()
		t := e.Type()
		for i := 0; i < e.NumField(); i++ {
			key := t.Field(i).Tag.Get("mapstructure")
			value := t.Field(i).Tag.Get("default")

			viper.SetDefault(key, value)
		}

		_ = viper.Unmarshal(cfg)
	}

	return cfg
}
