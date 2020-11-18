package nacos

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/ini"
	"github.com/tietang/props/v3/kvs"
	"github.com/tietang/props/v3/yam"
	"strings"
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type NacosClientConfigSource struct {
	NacosClientPropsConfigSource
}

func NewNacosClientConfigSource(address, group, dataId, namespaceId string) *NacosClientConfigSource {
	s := &NacosClientConfigSource{}

	name := strings.Join([]string{"Nacos", address}, ":")
	s.name = name
	s.DataId = dataId
	s.Group = group
	s.NamespaceId = namespaceId
	s.Values = make(map[string]string)
	s.NacosClientPropsConfigSource = *NewNacosClientPropsConfigSource(address, group, dataId, namespaceId)
	s.init()

	return s
}

func NewNacosClientCompositeConfigSource(address, group, tenant string, dataIds []string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "NacosKevValue"
	for _, dataId := range dataIds {
		c := NewNacosPropsConfigSource(address, group, dataId, tenant)
		s.Add(c)
	}

	return s
}

func (s *NacosClientConfigSource) init() {

	cr, err := s.get()
	if err != nil {
		log.Error(err)
		return
	}
	contentType := kvs.ContentType(cr.ContentType)
	if contentType == "text" {
		contentType = kvs.ReadContentType(cr.Content)
	}
	if contentType == kvs.ContentProps || contentType == kvs.ContentProperties {
		s.findProperties(cr.Content)
	} else if contentType == kvs.ContentIni {
		s.findIni(cr.Content)
	} else if contentType == kvs.ContentYaml || contentType == kvs.ContentYam || contentType == kvs.ContentYml {
		s.findYaml(cr.Content)
	} else {
		log.Warn("Unsupported format：", s.ContentType)
	}

}

func (s *NacosClientConfigSource) watchContext() {

}

func (s *NacosClientConfigSource) Close() {
}

func (s *NacosClientConfigSource) findYaml(content string) {
	props := yam.ByYaml(content)
	s.SetAll(props.Values)
}

func (s *NacosClientConfigSource) findIni(content string) {
	props := ini.ByIni(content)
	s.SetAll(props.Values)
}

func (s *NacosClientConfigSource) findProperties(content string) {
	props := kvs.ByProperties(content)
	s.SetAll(props.Values)
}

func (s *NacosClientConfigSource) registerProps(key, value string) {
	s.Set(strings.TrimSpace(key), strings.TrimSpace(value))

}

func (s *NacosClientConfigSource) Name() string {
	return s.name
}

func (h *NacosClientConfigSource) get() (cr *ConfigRes, err error) {

	cp := vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
	}
	//if len(h.AppName) > 0 {
	//	cp.AppName = h.AppName
	//}
	content, err := h.Client.GetConfig(cp)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cr = &ConfigRes{}
	err = json.Unmarshal([]byte(content), cr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return cr, err
}
