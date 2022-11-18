package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/ini"
	"github.com/tietang/props/v3/kvs"
	"github.com/tietang/props/v3/yam"
	"strings"
)

// 通过key/value来组织，过滤root prefix后，替换/为.作为properties key
// 配置内容开头是否包含一下开头：#@, ;@, //@, @，并定义配置内容格式的信息，
// 比如：;@ini , #@yaml, #@yml等.
// 支持的格式有：;@ini，#@yaml, #@yml, #@yam，#@props，#@properties，
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
	s.NacosClientPropsConfigSource = *newNacosClientPropsConfigSource(address, group, dataId, namespaceId)
	s.NacosClientPropsConfigSource.listenConfig()
	s.init()

	return s
}

func NewNacosClientCompositeConfigSource(address, group, namespaceId string, dataIds []string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "NacosKevValue"
	for _, dataId := range dataIds {
		c := NewNacosClientConfigSource(address, group, dataId, namespaceId)
		s.Add(c)
	}

	return s
}

func (s *NacosClientConfigSource) init() {

	content, err := s.get()
	if err != nil {
		log.Error(err)
		return
	}
	contentType := kvs.ReadContentType(content)
	s.ContentType = string(contentType)
	if contentType == kvs.ContentProps || contentType == kvs.ContentProperties {
		s.findProperties(content)
	} else if contentType == kvs.ContentIni {
		s.findIni(content)
	} else if contentType == kvs.ContentYaml || contentType == kvs.ContentYam || contentType == kvs.ContentYml {
		s.findYaml(content)
	} else {
		log.Warn("[Nacos] Unsupported config format：", contentType,
			" ,请检查配置内容开头是否包含一下开头：#@, ;@, //@, @，并定义配置内容格式的信息，比如：;@ini , #@yaml, #@yml等.支持的格式有：;@ini，#@yaml, #@yml, #@yam，#@props，#@properties，")
	}

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

func (h *NacosClientConfigSource) get() (cr string, err error) {

	cp := vo.ConfigParam{
		DataId: h.DataId,
		Group:  h.Group,
	}
	//if len(h.AppName) > 0 {
	//	cp.AppName = h.AppName
	//}
	return h.Client.GetConfig(cp)
	//if err != nil {
	//	log.Error(err)
	//	return "", err
	//}
	//

	//cr = &ConfigRes{}
	//err = json.Unmarshal([]byte(content), cr)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//return cr, err
}
