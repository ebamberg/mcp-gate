package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ebamberg/mcp-gate/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AppConfig struct {
	App struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"app"`
	Namespace string `mapstructure:"namespace"`
}

var Config AppConfig

func initConfig() {
	// Set up viper to read the config.yaml file
	viper.SetConfigName("config") // Config file name without extension
	viper.SetConfigType("yaml")   // Config file type
	viper.AddConfigPath(".")      // Look for the config file in the current directory

	// set default values for the config
	viper.SetDefault("app.name", "mcp-gate") // Set a default value for app.name

	/*
	   AutomaticEnv will check for an environment variable any time a viper.Get request is made.
	   It will apply the following rules.
	       It will check for an environment variable with a name matching the key uppercased and prefixed with the EnvPrefix if set.
	*/
	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")                              // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // this is useful e.g. want to use . in Get() calls, but environmental variables to use _ delimiters (e.g. app.port -> APP_PORT)

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Can't read config:", err)
	} else {
		// Unmarshal the config file into the AppConfig struct
		err = viper.Unmarshal(&Config)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Can't read config:", err)
			os.Exit(1)
		}
	}
}

func configLogging() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("MCP-GATE: ")
}

func main() {

	configLogging()
	cobra.OnInitialize(initConfig)
	cmd.Execute()

}
