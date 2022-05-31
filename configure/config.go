package configure

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type ServerCfg struct {
	ConfigFile  string `mapstructure:"config_file" json:"config_file"`
	Port        string `mapstructure:"port" json:"port"`
	Environment string `mapstructure:"environment" json:"environment"`

	TwitchClientId     string `mapstructure:"twitch_client_id" json:"twitch_client_id"`
	TwitchClientSecret string `mapstructure:"twitch_client_secret" json:"twitch_client_secret"`
	TwitchRedirectURI  string `mapstructure:"twitch_redirect_uri" json:"twitch_redirect_uri"`

	TwitcEventSubSecret string `mapstructure:"twitch_eventsub_secret" json:"twitch_eventsub_secret"`

	MongoURI string `mapstructure:"mongo_uri" json:"mongo_uri"`
}

var defaultConfig = ServerCfg{
	ConfigFile: "config.yaml",
}

var Config = viper.New()

func init() {
	Config.SetConfigFile(Config.GetString("config_file"))
	Config.SetConfigType("yaml")
	Config.AddConfigPath("./")
	Config.AddConfigPath("/go/src/github.com/redis_docker")

	err := Config.ReadInConfig()
	if err != nil {
		fmt.Println("Fatal error config file: default \n", err)
		os.Exit(1)
	}

	// Environment
	Config.AutomaticEnv()
	Config.SetEnvPrefix("GGAPI")
	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	Config.AllowEmptyEnv(true)
}
