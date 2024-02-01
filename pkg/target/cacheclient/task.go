package cacheclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/skhatri/go-fns/lib/fs"
	"github.com/skhatri/go-http-cache/pkg/conf"
	"github.com/skhatri/go-http-cache/pkg/target/model"
	"github.com/skhatri/go-logger/logging"
	"io"
	"os"
	"time"
)

var logger = logging.NewLogger("fs")

type CacheData struct {
	StatusCode int                 `json:"statusCode"`
	Data       []byte              `json:"data"`
	Headers    map[string][]string `json:"headers"`
}

func (cd *CacheData) ToResponse() *model.Response {
	return &model.Response{
		StatusCode: cd.StatusCode,
		Headers:    cd.Headers,
		Data:       io.NopCloser(bytes.NewReader(cd.Data)),
	}
}

type _CacheClient struct {
	cache Cache
}

func (fo *_CacheClient) Invoke(req model.Request) (*model.Response, error) {
	if res, ok := fo.cache.Get(req.Key()); ok {
		logger.WithTask("http-fs-cache").WithAttribute("key", req.Key()).WithMessage("served from cache").Debug()
		return res.ToResponse(), nil
	}
	return nil, fmt.Errorf("not available")
}

func (fo *_CacheClient) OnNotify(req model.Request, res *model.Response) {
	data, _ := io.ReadAll(res.Data)
	fo.cache.Store(req.Key(), CacheData{
		StatusCode: res.StatusCode,
		Headers:    res.Headers,
		Data:       data,
	})
}

type Cache interface {
	Name() string
	Get(string) (*CacheData, bool)
	Store(string, CacheData) bool
}

type mapCache struct {
	cache map[string]CacheData
}

func (mc *mapCache) Get(key string) (*CacheData, bool) {
	if data, ok := mc.cache[key]; ok {
		return &data, true
	}
	return nil, false
}
func (mc *mapCache) Store(key string, data CacheData) bool {
	mc.cache[key] = data
	return true
}

func (mc *mapCache) Name() string {
	return "map"
}

type fileCache struct {
	target string
}

func (fc *fileCache) Get(key string) (*CacheData, bool) {
	dataFile := fmt.Sprintf("%s/%s/%s", fc.target, key, "data.json")
	if _, err := os.Stat(dataFile); err != nil {
		return nil, false
	}

	fileRdr, rErr := os.Open(dataFile)
	if rErr != nil {
		return nil, false
	}
	cacheData := CacheData{}
	jErr := json.NewDecoder(fileRdr).Decode(&cacheData)
	if jErr != nil {
		return nil, false
	}
	if cacheData.Headers == nil {
		cacheData.Headers = make(map[string][]string)
	}
	cacheData.Headers["Date"] = []string{
		time.Now().Format(time.RFC1123),
	}
	cacheData.Headers["Cache-Hit"] = []string{"true"}
	cacheData.Headers["Cache-Engine"] = []string{fc.Name()}
	return &cacheData, true
}
func (fc *fileCache) Store(key string, data CacheData) bool {
	dataDir := fmt.Sprintf("%s/%s", fc.target, key)
	createErr := fs.CreateDirIfNotExists(dataDir)
	if createErr != nil {
		return false
	}
	targetFile := fmt.Sprintf("%s/%s", dataDir, "data.json")
	fw, err := os.Create(targetFile)
	if err != nil {
		return false
	}
	err = json.NewEncoder(fw).Encode(data)
	if err != nil {
		return false
	}
	return true
}
func (fc *fileCache) Name() string {
	return "fs"
}

func New(cacheSettings conf.Cache) model.ResourceClient {
	var cacheEngine Cache
	switch cacheSettings.Engine {
	case "fs":
		cacheEngine = &fileCache{
			cacheSettings.Location,
		}
	default:
		cacheEngine = &mapCache{
			cache: make(map[string]CacheData),
		}
	}
	return &_CacheClient{
		cache: cacheEngine,
	}
}
