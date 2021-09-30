package config

import (
	"encoding/json"
	"strings"

	"github.com/Strum355/log"
	"github.com/spf13/viper"
)

// Load fetches the config from environment variables or uses the defined defaults
func Load() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	loadDefaults()
	viper.AutomaticEnv()
}

func loadDefaults() {
	viper.SetDefault("alfred.commport", 7070)
	viper.SetDefault("alfred.token", "")
	viper.SetDefault("alfred.apptoken", "")
}

func PrintSettings() {
	settings := viper.AllSettings()
	// settings["alfred"].(map[string]interface{})["token"] = "[token]"
	// settings["alfred"].(map[string]interface{})["apptoken"] = "[apptoken]"

	out, _ := json.MarshalIndent(settings, "", "\t")
	log.Debug("config:\n" + string(out))
}
