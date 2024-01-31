package target

import (
	"fmt"
	"github.com/skhatri/go-http-cache/pkg/conf"
	"github.com/skhatri/go-http-cache/pkg/target/filesystem"
	"github.com/skhatri/go-http-cache/pkg/target/httpcall"
	"github.com/skhatri/go-http-cache/pkg/target/model"
)

type wrapper struct {
	clients []model.ResourceClient
}

func (wr *wrapper) Invoke(req model.Request) (*model.Response, error) {
	for _, c := range wr.clients {
		res, _ := c.Invoke(req)
		if res != nil {
			return res, nil
		}
	}
	return nil, fmt.Errorf("no request handler")
}

func NewResourceClient(cacheSettings conf.Cache) model.ResourceClient {
	fsclient := filesystem.New(cacheSettings)
	delegates := []model.ResourceClient{
		fsclient,
		httpcall.New(fsclient.(model.Notifier)),
	}
	return &wrapper{
		clients: delegates,
	}
}
