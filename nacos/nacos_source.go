package nacos

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"github.com/tietang/props/yam"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
)

const (
	ENDPOINT_GETALL_REQUEST = ENDPOINT_GET_REQUEST + "&show=all"
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type NacosConfigSource struct {
	NacosPropsConfigSource
}

func NewNacosConfigSource(address, group, dataId, tenant string) *NacosConfigSource {
	s := &NacosConfigSource{}
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

func NewNacosCompositeConfigSource(address, group, tenant string, dataIds []string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "NacosKevValue"
	for _, dataId := range dataIds {
		c := NewNacosPropsConfigSource(address, group, dataId, tenant)
		s.Add(c)
	}

	return s
}

func (s *NacosConfigSource) init() {

	cr, err := s.get()
	if err != nil {
		log.Error(err)
		return
	}
	if cr.ContentType == "properties" {
		s.findProperties(cr.Content)
	} else if cr.ContentType == "ini" {
		s.findIni(cr.Content)
	} else if cr.ContentType == "yaml" {
		s.findYaml(cr.Content)
	} else {
		log.Warn("Unsupported format：", cr.ContentType)
	}

}

func (s *NacosConfigSource) watchContext() {

}

func (s *NacosConfigSource) Close() {
}

func (s *NacosConfigSource) findYaml(content string) {
	props := yam.ByYaml(content)
	s.SetAll(props.Values)
}

func (s *NacosConfigSource) findIni(content string) {
	props := ini.ByIni(content)
	s.SetAll(props.Values)
}

func (s *NacosConfigSource) findProperties(content string) {
	props := kvs.ByProperties(content)
	s.SetAll(props.Values)
	//
	//sep := s.LineSeparator
	//if sep == "" {
	//	sep = NACOS_LINE_SEPARATOR
	//}
	//kvsep := s.KVSeparator
	//if kvsep == "" {
	//	kvsep = NACOS_KV_SEPARATOR
	//}
	//lines := strings.Split(content, sep)
	//
	//for _, l := range lines {
	//
	//	i := strings.Index(l, kvsep)
	//	if i <= 0 {
	//		continue
	//	}
	//	key := string(l[:i])
	//	value := string(l[i+1:])
	//	s.registerProps(key, value)
	//	//log.Info(key,"=",value)
	//}

}

func (s *NacosConfigSource) registerProps(key, value string) {
	s.Set(strings.TrimSpace(key), strings.TrimSpace(value))

}

func (s *NacosConfigSource) Name() string {
	return s.name
}

func (h *NacosConfigSource) Next() string {

	nv := atomic.AddUint32(&h.lastCt, 1)
	size := len(h.servers)
	if size == 0 {
		panic(errors.New("not found server."))
	}
	index := int(nv) % size
	selected := h.servers[index]
	return selected
}

func (h *NacosConfigSource) get() (cr *ConfigRes, err error) {
	base := h.Next()
	//?dataId=%s&group=%s&tenant=%s&show=all&
	url := fmt.Sprintf(ENDPOINT_GETALL_REQUEST, base, h.DataId, h.Group, h.Tenant)

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
	cr = &ConfigRes{}
	err = json.Unmarshal(respBody, cr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return cr, err
}
