package cmd

import (
	"fmt"
	"net"
	"strings"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/spf13/viper"
)

var App *EnvConfig

type EnvConfig struct {
	Environment   string
	RedisHost     string
	RedisPort     int
	RedisUsername string
	RedisPassword string
	DatabaseUrl   string
	GitHubToken   string
}

func isValidHost(value any) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}
	if strings.ToLower(s) == "localhost" {
		return nil
	}
	if ip := net.ParseIP(s); ip == nil {
		return nil
	}
	if err := is.URL.Validate(s); err == nil {
		return nil
	}
	return fmt.Errorf("must be 'localhost' or a valid URL/IP address")
}

func (e *EnvConfig) Validate() error {
	return v.ValidateStruct(e,
		v.Field(&e.Environment, v.Required, v.In("development", "production")),
		v.Field(&e.RedisHost, v.Required, v.By(isValidHost)),
		v.Field(&e.RedisPort, v.Required),
		v.Field(&e.RedisUsername),
		v.Field(&e.RedisPassword),
		v.Field(&e.DatabaseUrl, v.Required, is.URL),
		v.Field(&e.GitHubToken, v.Required),
	)
}

func SetupEnv() error {

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	App = &EnvConfig{
		Environment:   viper.GetString("environment"),
		RedisHost:     viper.GetString("redis.host"),
		RedisPort:     viper.GetInt("redis.port"),
		RedisUsername: viper.GetString("redis.username"),
		RedisPassword: viper.GetString("redis.password"),
		DatabaseUrl:   viper.GetString("database.url"),
		GitHubToken:   viper.GetString("github.personal_access_token"),
	}
	if err := App.Validate(); err != nil {
		return err
	}
	return nil
}
