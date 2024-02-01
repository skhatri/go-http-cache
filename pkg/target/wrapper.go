package target

import (
	"fmt"
	"github.com/skhatri/go-http-cache/pkg/conf"
	"github.com/skhatri/go-http-cache/pkg/target/cacheclient"
	"github.com/skhatri/go-http-cache/pkg/target/httpcall"
	"github.com/skhatri/go-http-cache/pkg/target/model"
	"strings"
)

type wrapper struct {
	clients []model.ResourceClient
}

func (wr *wrapper) Invoke(req model.Request) (*model.Response, error) {
	errs := make([]string, 0)
	for _, c := range wr.clients {
		res, err := c.Invoke(req)
		if err != nil {
			errs = append(errs, err.Error())
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, fmt.Errorf("no request handler. delegate errors [%s]", strings.Join(errs, ", "))
}

func NewResourceClient(cacheSettings conf.Cache) model.ResourceClient {
	cacheClient := cacheclient.New(cacheSettings)
	delegates := []model.ResourceClient{
		cacheClient,
		httpcall.New(cacheClient.(model.Notifier)),
	}
	return &wrapper{
		clients: delegates,
	}
}
