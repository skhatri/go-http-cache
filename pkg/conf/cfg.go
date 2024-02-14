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
	Engine   string       `yaml:"engine,omitempty"`
	Location string       `yaml:"location,omitempty"`
	Options  CacheOptions `yaml:"options,omitempty"`
}
type CacheOptions struct {
	LogMiss           *bool `yaml:"log-miss,omitempty"`
	LogHit            *bool `yaml:"log-hit,omitempty"`
	IgnoreHeaders     *bool `yaml:"ignore-headers,omitempty"`
	LogRequestHeaders *bool `yaml:"log-request-headers,omitempty"`
}

func (co *CacheOptions) ShouldLogMiss() bool {
	return co.LogMiss != nil && *co.LogMiss
}

func (co *CacheOptions) ShouldLogHit() bool {
	return co.LogHit != nil && *co.LogHit
}

func (co *CacheOptions) ShouldIgnoreHeaders() bool {
	return co.IgnoreHeaders != nil && *co.IgnoreHeaders
}

func (co *CacheOptions) ShouldLogRequestHeaders() bool {
	return co.LogRequestHeaders != nil && *co.LogRequestHeaders
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
	if ignoreHeaderOverride := os.Getenv("IGNORE_HEADERS"); ignoreHeaderOverride != "" {
		ignoreFlag := strings.ToLower(ignoreHeaderOverride) == "true"
		Configuration.Cache.Options.IgnoreHeaders = &ignoreFlag
	}

	if logMiss := os.Getenv("LOG_MISS"); logMiss != "" {
		missFlag := strings.ToLower(logMiss) == "true"
		Configuration.Cache.Options.LogMiss = &missFlag
	}
	if logRequestHeaders := os.Getenv("LOG_REQUEST_HEADERS"); logRequestHeaders != "" {
		requestHeadersFlag := strings.ToLower(logRequestHeaders) == "true"
		Configuration.Cache.Options.LogRequestHeaders = &requestHeadersFlag
	}

	if err != nil {
		panic(err)
	}
}
