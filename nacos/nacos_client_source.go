package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
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
	kvs.MapProperties
	name string
	// Required Configuration ID. Use a naming rule similar to package.class (for jiebaexample, com.taobao.tc.refund.log.level) to ensure global uniqueness. It is recommended to indicate business meaning of the configuration in the "class" section. Use lower case for all characters. Use alphabetical letters and these four special characters (".", ":", "-", "_") only. Up to 256 characters are allowed.
	DataId string
	// Required Configuration group. To ensure uniqueness, format such as product name: module name (for jiebaexample, Nacos:Test) is preferred. Use alphabetical letters and these four special characters (".", ":", "-", "_") only. Up to 128 characters are allowed.
	Group string

	//Tenant information. It corresponds to the Namespace field in Nacos.
	//Tenant      string
	NamespaceId   string
	ContentType   string
	AppName       string
	ClientConfig  *constant.ClientConfig
	ServerConfigs []constant.ServerConfig
	Client        config_client.IConfigClient
	IsListening   bool
}

func NewNacosClientConfigSource(address, group, namespaceId, dataId string) *NacosClientConfigSource {
	s := &NacosClientConfigSource{}
	name := strings.Join([]string{"Nacos", address}, ":")
	s.name = name
	s.DataId = dataId
	s.Group = group
	s.NamespaceId = namespaceId
	s.IsListening = true
	s.Values = make(map[string]string)
	var err error
	s.ClientConfig, s.ServerConfigs, s.Client, err = buildNacos(address, namespaceId)
	if err != nil {
		log.Panic("error create ConfigClient: ", err)
	}
	s.Values = make(map[string]string)

	s.init()
	return s
}

func NewNacosClientCompositeConfigSource(address, group, namespaceId string, dataIds []string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "NacosCompositeKevValue"
	for _, dataId := range dataIds {
		c := NewNacosClientConfigSource(address, group, namespaceId, dataId)
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
	s.parseConfig(content)
	if s.IsListening {
		s.listenConfig()
	}
}

func (s *NacosClientConfigSource) parseConfig(content string) {

	contentType := kvs.GetContentTypeByName(s.DataId)
	if contentType == kvs.ContentUnknown {
		contentType = kvs.ReadContentType(content)
	}
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
func (s *NacosClientConfigSource) CancelListening() {
	cp := vo.ConfigParam{
		DataId: s.DataId,
		Group:  s.Group,
	}
	err := s.Client.CancelListenConfig(cp)
	if err != nil {
		log.Error(err)
	}
	s.IsListening = false
}

func (s *NacosClientConfigSource) listenConfig() {
	cp := vo.ConfigParam{
		DataId: s.DataId,
		Group:  s.Group,
	}
	//if len(s.AppName) > 0 {
	//	cp.AppName = s.AppName
	//}
	cp.OnChange = func(namespace, group, dataId, data string) {
		if s.IsListening {
			s.parseConfig(data)
			log.Infof("changed config: %s %s %s", namespace, group, dataId)
		}
	}
	err := s.Client.ListenConfig(cp)
	if err != nil {
		log.Error("listen config： ", err)
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
