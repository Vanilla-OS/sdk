package conf

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/vanilla-os/sdk/pkg/v1/conf/types"
)

// InitConfig initializes a viper configuration and returns a pointer to the
// configuration struct.
//
// Example:
//
//	confStruct := conf.ConfigStruct{
//		Place:    "Gotham",
//		Event:    "Joker's Robbery",
//		Duration: 24,
//	}
//
//	opts := types.ConfigOptions{
//		Domain: "org.gotham.events",
//		Prefix: "/tmp",
//		Type:   "yml",
//	}
//
//	config, err := conf.InitConfig[conf.ConfigStruct](opts)
//	if err != nil {
//		fmt.Printf("error initializing config: %v", err)
//	}
//
//	fmt.Printf("The event %s at %s will last for %d hours", config.Event, config.Place, config.Duration)
func InitConfig[T any](opts types.ConfigOptions) (*T, error) {
	var config T

	configPaths := buildConfigPaths(opts.Prefix, opts.Domain, opts.Type)
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config")
	if opts.Type != "" {
		viper.SetConfigType(opts.Type)
	} else {
		viper.SetConfigType("yaml")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// buildConfigPaths returns a list of paths where the configuration file might
// be located, in order of priority.
func buildConfigPaths(prefix, domain, fileType string) []string {
	user, err := user.Current()
	if err != nil {
		return []string{}
	}
	userDir := user.HomeDir

	configPaths := []string{
		filepath.Join(".", "conf", domain),
		filepath.Join(prefix, userDir, domain),
		filepath.Join(prefix, "/etc", domain),
		filepath.Join(prefix, "/usr/share", domain),
	}

	return configPaths
}
