package nacos

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/nacos_error"
	"github.com/nacos-group/nacos-sdk-go/common/util"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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

	CONFIG_BASE_PATH   = "/v1/cs"
	CONFIG_PATH        = CONFIG_BASE_PATH + "/configs"
	CONFIG_LISTEN_PATH = CONFIG_BASE_PATH + "/configs/listener"
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
	Tenant         string
	ContentType    string
	AppName        string
	NamespaceId    string
	TimeoutMs      uint64 //10 * 1000, //http请求超时时间，单位毫秒
	ListenInterval uint64 //30 * 1000, //监听间隔时间，单位毫秒（仅在ConfigClient中有效）
	BeatInterval   uint64 //5 * 1000, //心跳间隔时间，单位毫秒（仅在ServiceClient中有效）

	LineSeparator string
	KVSeparator   string
	//
	name     string
	servers  []string
	lastCt   uint32
	mutex    sync.Mutex
	OnChange func(namespace, group, dataId, data string)
	lastMD5  string
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
	s.watchContext()
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
	return h.getInner(h.DataId, h.Group)
}

func (h *NacosPropsConfigSource) getInner(dataId, group string) (body []byte, err error) {
	base := h.Next()
	//?dataId=%s&group=%s&tenant=%s
	//show=all&
	url := fmt.Sprintf(ENDPOINT_GET_REQUEST, base, dataId, group, h.Tenant)

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
func (s *NacosPropsConfigSource) watchContext() {
	s.ListenConfig()
}

func (cp *NacosPropsConfigSource) ListenConfig() {
	go func() {
		for {
			params := make(map[string]string, 0)
			var listeningConfigs string
			md5 := cp.lastMD5

			if len(cp.Tenant) > 0 {
				listeningConfigs += cp.DataId + constant.SPLIT_CONFIG_INNER + cp.Group + constant.SPLIT_CONFIG_INNER +
					md5 + constant.SPLIT_CONFIG_INNER + cp.Tenant + constant.SPLIT_CONFIG
			} else {
				listeningConfigs += cp.DataId + constant.SPLIT_CONFIG_INNER + cp.Group + constant.SPLIT_CONFIG_INNER +
					md5 + constant.SPLIT_CONFIG
			}

			params["Listening-Configs"] = listeningConfigs

			rs, err := cp.RequestListenConfig(params, cp.Tenant, "", "")
			if err != nil {
				continue
			}
			cp.updateLocalConfig(rs)
		}
	}()
}

func (n *NacosPropsConfigSource) updateLocalConfig(changed string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	changedConfigs := strings.Split(changed, "%01")
	for _, config := range changedConfigs {
		attrs := strings.Split(config, "%02")
		if len(attrs) == 2 {
			data, err := n.getInner(attrs[0], attrs[1])
			if err != nil {
				log.Println("[client.updateLocalConfig] update config failed:", err.Error())
			} else {
				//n.putLocalConfig(vo.ConfigParam{
				//	DataId:  attrs[0],
				//	Group:   attrs[1],
				//	Content: content,
				//})
				content := string(data)
				// call listener:
				//decrept, _ := n.decrypt(attrs[0], content)
				n.OnChange("", attrs[1], attrs[0], content)
				if len(content) > 0 {
					n.lastMD5 = util.Md5(content)
				}
			}
		} else if len(attrs) == 3 {
			data, err := n.getInner(attrs[0], attrs[1])
			if err != nil {
				log.Println("[client.updateLocalConfig] update config failed:", err.Error())
			} else {
				//client.putLocalConfig(vo.ConfigParam{
				//	DataId:  attrs[0],
				//	Group:   attrs[1],
				//	Content: content,
				//})
				content := string(data)
				// call listener:
				//decrept, _ := client.decrypt(attrs[0], content)
				n.OnChange(attrs[2], attrs[1], attrs[0], string(data))
				if len(content) > 0 {
					n.lastMD5 = util.Md5(content)
				}
			}

		}
	}
	log.Println("[client.updateLocalConfig] update config complete")
	//log.Println("[client.localConfig] ", client.localConfigs)
}

func (cp *NacosPropsConfigSource) RequestListenConfig(params map[string]string, tenant, accessKey, secretKey string) (result string, err error) {

	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	header.Add("Long-Pulling-Timeout", strconv.FormatUint(cp.ListenInterval, 10))
	header.Add("accessKey", accessKey)
	header.Add("secretKey", secretKey)

	header.Add("Client-Version", constant.CLIENT_VERSION)
	header.Add("User-Agent", constant.CLIENT_VERSION)
	//header.Add("Accept-Encoding","gzip,deflate,sdch"}
	header.Add("Connection", "Keep-Alive")
	header.Add("exConfigInfo", "true")
	uid, _ := uuid.NewV4()

	header.Add("RequestId", uid.String())
	header.Add("Request-Module", "Naming")
	header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	header.Add("Spas-AccessKey", accessKey)
	timeStamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	header.Add("Timestamp", timeStamp)

	resource := ""

	if len(params["tenant"]) != 0 {
		resource = params["tenant"] + "+" + params["group"]
	} else {
		resource = params["group"]
	}
	signature := ""
	if resource == "" {
		signature = signWithhmacSHA1Encrypt(timeStamp, secretKey)
	} else {
		signature = signWithhmacSHA1Encrypt(resource+"+"+timeStamp, secretKey)
	}
	header.Add("Spas-Signature", signature)
	log.Printf("[client.ListenConfig] request params:%+v header:%+v \n", params, header)

	//使用cp.ListenInterval代替超时。
	res, err := post(constant.CONFIG_LISTEN_PATH, header, cp.ListenInterval, params)
	if err != nil {
		return
	}
	var bytes []byte
	bytes, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}
	result = string(bytes)
	if res.StatusCode == 200 {
		return
	} else {
		err = nacos_error.NewNacosError(strconv.Itoa(res.StatusCode), string(bytes), nil)
		return
	}

	return result, err
}

func signWithhmacSHA1Encrypt(encryptText, encryptKey string) string {
	//hmac ,use sha1
	key := []byte(encryptKey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(encryptText))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func post(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error) {
	client := http.Client{}
	client.Timeout = time.Millisecond * time.Duration(timeoutMs)
	var body string
	for key, value := range params {
		if len(value) > 0 {
			body += key + "=" + value + "&"
		}
	}
	if strings.HasSuffix(body, "&") {
		body = body[:len(body)-1]
	}
	log.Info(path)
	log.Info(body)
	request, errNew := http.NewRequest(http.MethodPost, path, strings.NewReader(body))
	if errNew != nil {
		err = errNew
		return
	}
	request.Header = header
	resp, errDo := client.Do(request)
	if errDo != nil {
		err = errDo
	} else {
		response = resp
	}
	return
}
