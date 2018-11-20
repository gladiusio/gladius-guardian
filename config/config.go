package config

import (
	"fmt"
	"strings"

	"github.com/gladiusio/gladius-common/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SetupConfig - Setup a config file and add some default values
func SetupConfig(configFilePath string) {
	viper.SetConfigName("gladius-guardian")
	viper.AddConfigPath(configFilePath)

	// Setup env variable handling
	viper.SetEnvPrefix("GUARDIAN")
	r := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(r)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Warn(fmt.Errorf("error reading config file: %s, using defaults", err))
	}

	// add some defaults
	AddViperDefault("NetworkdExecutable", "gladius-networkd")
	AddViperDefault("ControldExecutable", "gladius-controld")
	AddViperDefault("Ports.Guardian", 7791)
	AddViperDefault("Ports.EdgeD", 7946)
	AddViperDefault("Ports.NetworkGateway", 3001)

	// Get gladius base for the various services
	base, err := utils.GetGladiusBase()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Warn("Couldn't get Gladius base")
	}

	// Add a default environment so that we can set the gladius base of our sub
	// processes
	AddViperDefault("DefaultEnvironment", []string{"GLADIUSBASE=" + base})
	AddViperDefault("MaxLogLines", 1000) // Max number of log lines to keep in ram for each service

	// Setup logging level
	switch loglevel := viper.GetString("LogLevel"); loglevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

// AddViperDefault - add a default key
func AddViperDefault(key string, defaultValue interface{}) string {
	viper.SetDefault(key, defaultValue)

	return key
}
