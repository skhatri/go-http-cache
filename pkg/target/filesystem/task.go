package filesystem

import (
	"bytes"
	"fmt"
	"github.com/skhatri/go-http-cache/pkg/conf"
	"github.com/skhatri/go-http-cache/pkg/target/model"
	"github.com/skhatri/go-logger/logging"
	"io"
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

type _fsOp struct {
	target string
	cache  map[string]CacheData
}

func (fo *_fsOp) Invoke(req model.Request) (*model.Response, error) {
	if res, ok := fo.cache[req.Key()]; ok {
		logger.WithTask("http-fs-cache").WithAttribute("key", req.Key()).WithMessage("served from cache").Debug()
		return res.ToResponse(), nil
	}
	return nil, fmt.Errorf("not available")
}

func (fo *_fsOp) OnNotify(req model.Request, res *model.Response) {
	data, _ := io.ReadAll(res.Data)
	fo.cache[req.Key()] = CacheData{
		StatusCode: res.StatusCode,
		Headers:    res.Headers,
		Data:       data,
	}
}
func New(cacheSettings conf.Cache) model.ResourceClient {
	return &_fsOp{
		target: cacheSettings.Location,
		cache:  make(map[string]CacheData),
	}
}
