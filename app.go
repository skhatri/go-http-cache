package main

import (
        "github.com/gin-gonic/gin"
        "github.com/skhatri/go-logger/logging"
        "github.com/skhatri/go-fns/lib/converters"
        "net/http"
        "os"
)


type Server struct {
  Address string `yaml:"address,omitempty"`
}

type Cache struct {
  Engine string `yaml:"engine,omitempty"`
  Location string `yaml:"location,omitempty"` 
}

type Config struct {
  Server Server `yaml:"server,omitempty"`
  Target []string `yaml:"target,omitempty"`
  Cache Cache `yaml:"location,omitempty"`  
}

func Configure() {
        configFile := "config.yaml"
        if cfg := os.Getenv("CONFIG_FILE"); cfg != "" {
          configFile = cfg
        }

        cf := &Config{}
      	err := converters.UnmarshalFile(configFile, cf)
        if err != nil {
          panic(err)
        } 
        gin.SetMode(gin.ReleaseMode)
        r := gin.Default()

        r.GET("/readiness", statusOk)
        r.GET("/liveness", statusOk)

        r.Run(cf.Server.Address)
}
var logger = logging.NewLogger("configure")

func statusOk(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
                "status": "OK",
        })
}

func main() {
   Configure()
}
