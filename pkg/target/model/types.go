package model

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
	"sort"
	"strings"
)

type Request struct {
	Method  string
	Url     string
	Headers map[string][]string
	Body    []byte
}

func (req *Request) Key(ignoreHeaders bool) string {
	hashObj := NewHashing().Write(req.Method).Write(req.Url)
	headersSerial := bytes.Buffer{}
	if !ignoreHeaders {
		keys := make([]string, 0)
		for k, _ := range req.Headers {
			key := strings.ToLower(k)
			if key == "user-agent" || 
			strings.HasPrefix(key, "authorization") ||
			strings.HasPrefix(key, "sec-") || 
			strings.HasPrefix(key, "accept-") {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			values := req.Headers[key]
			sort.Strings(values)
			headersSerial.WriteString(key)
			headersSerial.WriteString("=")
			headersSerial.WriteString(strings.Join(values, ","))
			headersSerial.WriteString(";")
		}
		hashObj.WriteBytes(headersSerial.Bytes())
	}
	if req.Body != nil {
		hashObj.WriteBytes(req.Body)
	}
	return hashObj.Sum()
}

func NewHashing() *Hashing {
	return &Hashing{
		ignoreCase: true,
		_hash:      md5.New(),
	}
}

type Hashing struct {
	ignoreCase bool
	_hash      hash.Hash
}

func (hs *Hashing) Write(s string) *Hashing {
	value := s
	if hs.ignoreCase {
		value = strings.ToLower(s)
	}
	hs._hash.Write([]byte(value))
	return hs
}

func (hs *Hashing) WriteBytes(b []byte) *Hashing {
	if b == nil {
		return hs
	}
	hs._hash.Write(b)
	return hs
}
func (hs *Hashing) Sum() string {
	data := hs._hash.Sum(nil)
	return hex.EncodeToString(data)
}

type Response struct {
	StatusCode     int
	Headers        map[string][]string
	Data           io.ReadCloser
	RequestHeaders map[string][]string
}
type ResourceClient interface {
	Invoke(request Request) (*Response, error)
}

type Notifier interface {
	OnNotify(Request, *Response)
}
