package conf

import (
	"github.com/skhatri/go-fns/lib/converters"
	"os"
)

type Server struct {
	Address string `yaml:"address,omitempty"`
}

type Cache struct {
	Engine   string `yaml:"engine,omitempty"`
	Location string `yaml:"location,omitempty"`
}

type Config struct {
	Server Server   `yaml:"server,omitempty"`
	Target []string `yaml:"target,omitempty"`
	Cache  Cache    `yaml:"location,omitempty"`
}

var Configuration *Config

func init() {
	configFile := "config.yaml"
	if cfg := os.Getenv("CONFIG_FILE"); cfg != "" {
		configFile = cfg
	}
	Configuration = &Config{}
	err := converters.UnmarshalFile(configFile, Configuration)
	if err != nil {
		panic(err)
	}
}
