package nacos

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
)

const (
	NACOS_LINE_SEPARATOR = "\n"
	NACOS_KV_SEPARATOR   = "="
	REQUEST_TEMPLATE     = "http://%s"

	ENDPOINT_GET            = "/nacos/v1/cs/configs"
	ENDPOINT_GET_REQUEST    = REQUEST_TEMPLATE + ENDPOINT_GET + "?dataId=%s&group=%s&tenant=%s"
	ENDPOINT_LISTEN         = "/nacos/v1/cs/configs/listener"
	ENDPOINT_LISTEN_REQUEST = REQUEST_TEMPLATE + ENDPOINT_LISTEN
	ENDPOINT_LISTEN_BODY    = "?dataId=%s&group=%s&tenant=%s"
)

var _ kvs.ConfigSource = new(NacosPropsConfigSource)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type NacosPropsConfigSource struct {
	kvs.MapProperties
	// Required Configuration ID. Use a naming rule similar to package.class (for example, com.taobao.tc.refund.log.level) to ensure global uniqueness. It is recommended to indicate business meaning of the configuration in the "class" section. Use lower case for all characters. Use alphabetical letters and these four special characters (".", ":", "-", "_") only. Up to 256 characters are allowed.
	DataId string
	// Required Configuration group. To ensure uniqueness, format such as product name: module name (for example, Nacos:Test) is preferred. Use alphabetical letters and these four special characters (".", ":", "-", "_") only. Up to 128 characters are allowed.
	Group string

	//Tenant information. It corresponds to the Namespace field in Nacos.
	Tenant      string
	ContentType string
	AppName     string
	NamespaceId string

	LineSeparator string
	KVSeparator   string
	//
	name    string
	servers []string
	lastCt  uint32
}

func NewNacosPropsConfigSource(address, group, dataId, tenant string) *NacosPropsConfigSource {
	s := &NacosPropsConfigSource{}
	s.servers = strings.Split(address, ",")
	name := strings.Join([]string{"Nacos", address}, ":")
	s.name = name
	s.DataId = dataId
	s.Group = group
	s.Tenant = tenant
	s.Values = make(map[string]string)
	s.init()

	return s
}

func NewNacosPropsCompositeConfigSource(address, group, tenant string, dataIds []string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "NacosKevValue"
	for _, dataId := range dataIds {
		c := NewNacosPropsConfigSource(address, group, dataId, tenant)
		s.Add(c)
	}

	return s
}

func (s *NacosPropsConfigSource) init() {
	s.findProperties()
}

func (s *NacosPropsConfigSource) watchContext() {

}

func (s *NacosPropsConfigSource) Close() {
}

func (s *NacosPropsConfigSource) findProperties() {

	data, err := s.get()
	if err != nil {
		log.Error(err)
		return
	}
	sep := s.LineSeparator
	if sep == "" {
		sep = NACOS_LINE_SEPARATOR
	}
	kvsep := s.KVSeparator
	if kvsep == "" {
		kvsep = NACOS_KV_SEPARATOR
	}
	lines := bytes.Split(data, []byte(sep))

	for _, l := range lines {

		i := bytes.Index(l, []byte(kvsep))
		if i <= 0 {
			continue
		}
		key := string(l[:i])
		value := string(l[i+1:])
		s.registerProps(key, value)
		//log.Info(key,"=",value)
	}

}

func (s *NacosPropsConfigSource) registerProps(key, value string) {
	s.Set(strings.TrimSpace(key), strings.TrimSpace(value))

}

func (s *NacosPropsConfigSource) Name() string {
	return s.name
}

func (h *NacosPropsConfigSource) Next() string {

	nv := atomic.AddUint32(&h.lastCt, 1)
	size := len(h.servers)
	if size == 0 {
		panic(errors.New("not found server."))
	}
	index := int(nv) % size
	selected := h.servers[index]
	return selected
}

func (h *NacosPropsConfigSource) get() (body []byte, err error) {
	base := h.Next()
	//?dataId=%s&group=%s&tenant=%s
	//show=all&
	url := fmt.Sprintf(ENDPOINT_GET_REQUEST, base, h.DataId, h.Group, h.Tenant)

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
	return respBody, err
}
