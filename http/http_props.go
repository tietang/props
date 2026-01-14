package http

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"io"
	"net/http"
)

/*
*
request url: http://127.0.0.1:8080/v1/props?namespace=base
response:
{

	"pkey1":{
	    "p.key1":"value1",
	    "p.key2":"value2",
	    "p1.p2.key1":10,
	    "time.p2.key2":"10s",
	},

	"pkey2":{
	    "p.key1":"value1",
	    "p.key2":"value2",
	    "p1.p2.key1":10,
	    "time.p2.key2":"10s",
	}

}
*/
var _ kvs.ConfigSource = new(HttpPropsConfigSource)

type HttpPropsConfigSource struct {
	kvs.MapProperties
	name       string
	url        string
	namespaces []string
	username   string
	password   string
}

func NewHttpPropsConfigSource(name string, url string, namespaces []string) *HttpPropsConfigSource {
	s := &HttpPropsConfigSource{}
	s.name = name
	s.Values = make(map[string]string)
	s.url = url
	s.namespaces = namespaces
	s.init()
	return s
}

func (s *HttpPropsConfigSource) init() {
	s.findProperties()
}

func (s *HttpPropsConfigSource) findProperties() {
	for _, namespace := range s.namespaces {
		res, err := http.Get(s.url + "?namespace=" + namespace)
		if err != nil || res.StatusCode != http.StatusOK {
			continue
		}
		dst := bytes.NewBufferString("")
		io.Copy(dst, res.Body)
		m := make(map[string]map[string]string)
		err = json.Unmarshal(dst.Bytes(), &m)
		if err != nil {
			log.Error(err)
		}
		for pkey, keyValues := range m {
			for key, value := range keyValues {
				s.Set(pkey+"."+key, value)
			}
		}
	}
}
