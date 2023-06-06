package env

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresHost string `mapstructure:"POSTGRES_HOST" default:"localhost"`
	PostgresPort int    `mapstructure:"POSTGRES_PORT" default:"5432"`
	PostgresDB   string `mapstructure:"POSTGRES_DB" default:"dbo_test"`
	PostgresUser string `mapstructure:"POSTGRES_USER" default:"postgres"`
	PostgresPass string `mapstructure:"POSTGRES_PASSWORD" default:"962812"`
	AppHost      string `mapstructure:"APP_HOST" default:"localhost"`
	AppPort      int    `mapstructure:"APP_PORT" default:"5001"`
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
