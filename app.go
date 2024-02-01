package main

import (
	"fmt"
	"github.com/skhatri/go-http-cache/pkg/conf"
	"github.com/skhatri/go-http-cache/pkg/target"
	"github.com/skhatri/go-http-cache/pkg/target/model"
	"github.com/skhatri/go-logger/logging"
	"io"
	"net/http"
)

var logger = logging.NewLogger("configure")

func Configure() {
	var cf = conf.Configuration
	logger.WithTask("http-go-cache").WithMessage("running server on %s", cf.Server.Address).Info()
	err := http.ListenAndServe(cf.Server.Address, &handler{
		cf,
		target.NewResourceClient(cf.Cache),
	})
	if err != nil {
		panic(err)
	}
}

type handler struct {
	cf       *conf.Config
	instance model.ResourceClient
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var body []byte
	if r.Method != "GET" {
		rdr := r.Body
		defer rdr.Close()
		data, err := io.ReadAll(rdr)
		body = data
		if err != nil {
			writeError(w, err)
			return
		}
	}

	url := fmt.Sprintf("%s%s", h.cf.Target[0], r.RequestURI)
	clientRequest := model.Request{
		Method:  r.Method,
		Headers: r.Header,
		Body:    body,
		Url:     url,
	}
	res, invErr := h.instance.Invoke(clientRequest)
	if invErr != nil {
		writeError(w, invErr)
		return
	}
	defer res.Data.Close()

	for k, values := range res.Headers {
		if k == "Date" {
			continue
		}
		for _, value := range values {
			w.Header().Add(k, value)
		}
	}
	w.WriteHeader(res.StatusCode)
	_, cpError := io.Copy(w, res.Data)
	if cpError != nil {
		logger.WithTask("print-response").WithError(cpError)
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Add("content-type", "text/plain")
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func main() {
	Configure()
}
