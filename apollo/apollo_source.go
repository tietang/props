package apollo

import (
	"encoding/json"
	"fmt"
	"github.com/shima-park/agollo"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"github.com/tietang/props/yam"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type ApolloConfigSource struct {
	kvs.MapProperties
	Apollo      agollo.Agollo
	name        string
	address     string
	appId       string
	namespaces  []string
	cluster     string
	clientIP    string
	contentType kvs.ContentType
}

func NewApolloConfigSource(address, appId string, namespaces []string) *ApolloConfigSource {
	s := &ApolloConfigSource{}
	s.name = "apollo:" + address
	s.appId = appId
	s.address = address
	s.namespaces = namespaces
	s.contentType = kvs.KeyValueContentType
	s.clientIP, _ = utils.GetExternalIP()
	s.cluster = "default"
	s.Values = make(map[string]string)
	s.init()

	return s
}

func NewApolloCompositeConfigSource(url, appId string, namespaces []string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "ApolloKevValue"
	c := NewApolloConfigSource(url, appId, namespaces)
	s.Add(c)
	return s
}
func (s *ApolloConfigSource) init() {
	for _, ns := range s.namespaces {
		kvs, err := s.GetConfigsFromCache(s.address, s.appId, s.cluster, ns)
		if err != nil {
			log.Error(err)
			continue
		}
		for key, value := range kvs {
			s.initValue(ns, key, value)
		}
	}
}

func (s *ApolloConfigSource) initValue(namespace, key, value string) {
	contentType := s.getContentType(namespace)
	if contentType == kvs.KeyValueContentType {
		contentType = s.getContentType(key)
	}
	if contentType == kvs.KeyValueContentType {
		s.Set(key, value)
	} else if contentType == kvs.ContentProps || contentType == kvs.ContentProperties {
		s.findProperties(value)
	} else if contentType == kvs.ContentIni {
		s.findIni(value)
	} else if contentType == kvs.ContentYaml || contentType == kvs.ContentYam || contentType == kvs.ContentYml {
		s.findYaml(value)
	} else {
		s.Set(key, value)
	}

}

func (s *ApolloConfigSource) getContentType(key string) (ctype kvs.ContentType) {
	idx := strings.LastIndex(key, ".")
	if idx == -1 || idx == len(key)-1 {
		ctype = kvs.KeyValueContentType
	} else {
		ctype = kvs.ContentType(key[idx+1:])
	}
	return
}

func (s *ApolloConfigSource) watchContext() {

}

func (s *ApolloConfigSource) Close() {
}

func (s *ApolloConfigSource) findYaml(content string) {
	props := yam.ByYaml(content)
	s.SetAll(props.Values)
}

func (s *ApolloConfigSource) findIni(content string) {
	props := ini.ByIni(content)
	s.SetAll(props.Values)
}

func (s *ApolloConfigSource) findProperties(content string) {
	props := kvs.ByProperties(content)
	s.SetAll(props.Values)
}

func (s *ApolloConfigSource) registerProps(key, value string) {
	s.Set(strings.TrimSpace(key), strings.TrimSpace(value))

}

func (s *ApolloConfigSource) Name() string {
	return s.name
}

func (c *ApolloConfigSource) GetConfigsFromCache(address, appID, cluster, namespace string) (kvs KeyValue, err error) {
	url := fmt.Sprintf("http://%s/configfiles/json/%s/%s/%s?ip=%s",
		address,
		url.QueryEscape(appID),
		url.QueryEscape(cluster),
		url.QueryEscape(namespace),
		c.clientIP,
	)

	//调用请求
	res, err := http.Get(url)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 如果出错就不需要close，因此defer语句放在err处理逻辑后面
	defer res.Body.Close()
	//处理response,读取Response body
	respBody, err := ioutil.ReadAll(res.Body)

	//
	if err := res.Body.Close(); err != nil {
		log.Error(err)
	}
	kvs = KeyValue{}
	err = json.Unmarshal(respBody, &kvs)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return kvs, err

}

type KeyValue map[string]string

type ConfigRes struct {
	AppID      string   `json:"appId"`          // appId: "AppTest",
	Cluster    string   `json:"cluster"`        // cluster: "default",
	Namespace  string   `json:"namespaceName"`  // namespaceName: "TEST.Namespace1",
	KeyValue   KeyValue `json:"configurations"` // configurations: {Name: "Foo"},
	ReleaseKey string   `json:"releaseKey"`     // releaseKey: "20181017110222-5ce3b2da895720e8"
}
