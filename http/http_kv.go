package http

import (
    "github.com/tietang/props/kvs"
    "net/http"
    "bytes"
    "io"
    log "github.com/sirupsen/logrus"
    "encoding/json"
)

/**
    request url: http://127.0.0.1:8080/v1/props?namespace=base
response:
{
"p.key1":"value1",
"p.key2":"value2",
"p1.p2.key1":10,
"time.p2.key2":"10s",

}

 */
type HttpKeyValueConfigSource struct {
    kvs.MapProperties
    name       string
    url        string
    namespaces []string
    username   string
    password   string
}

func NewHttpKeyValueConfigSource(name string, url string, namespaces []string) *HttpKeyValueConfigSource {
    s := &HttpKeyValueConfigSource{}
    s.name = name
    s.Values = make(map[string]string)
    s.url = url
    s.namespaces = namespaces
    s.init()
    return s
}
func (s *HttpKeyValueConfigSource) init() {
    s.findProperties()
}

func (s *HttpKeyValueConfigSource) findProperties() {
    for _, namespace := range s.namespaces {
        res, err := http.Get(s.url + "?namespace=" + namespace)
        if err != nil || res.StatusCode != http.StatusOK {
            continue
        }
        dst := bytes.NewBufferString("")
        io.Copy(dst, res.Body)
        m := make(map[string]string)
        err = json.Unmarshal(dst.Bytes(), &m)
        if err != nil {
            log.Error(err)
        }
        s.SetAll(m)
    }
}
