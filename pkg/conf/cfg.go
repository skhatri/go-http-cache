package conf

import (
	"github.com/skhatri/go-fns/lib/converters"
	"os"
	"strings"
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
	Cache  Cache    `yaml:"cache,omitempty"`
}

var Configuration *Config

func init() {
	configFile := "config.yaml"
	if cfg := os.Getenv("CONFIG_FILE"); cfg != "" {
		configFile = cfg
	}
	Configuration = &Config{}
	err := converters.UnmarshalFile(configFile, Configuration)

	if targetOverride := os.Getenv("TARGET"); targetOverride != "" {
		Configuration.Target = strings.Split(targetOverride, ",")
	}
	if addressOverride := os.Getenv("LISTEN_ADDRESS"); addressOverride != "" {
		Configuration.Server.Address = addressOverride
	}

	if err != nil {
		panic(err)
	}
}
