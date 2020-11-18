package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func update(group, dataId, tenant string, size, len int) map[string]string {
	m := make(map[string]string)
	c := client()
	content := ""

	for i := 0; i < size; i++ {
		key := "key-" + strconv.Itoa(i)
		for j := 0; j < len; j++ {
			kk := key + "." + "x" + strconv.Itoa(j)
			val := "v-" + strconv.Itoa(i) + strconv.Itoa(j)
			k := strings.Replace(kk, "/", ".", -1)
			//fmt.Println(key, k, value)
			m[k] = val
			content += k
			content += "="
			content += val
			content += "\n"
		}
	}
	vo := vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: "",
	}
	b, e := c.PublishConfig(vo)
	fmt.Println(b, e)

	return m
}

func client() config_client.IConfigClient {
	clientConfig := constant.ClientConfig{
		TimeoutMs:            10 * 1000,            //http请求超时时间，单位毫秒
		ListenInterval:       30 * 1000,            //监听间隔时间，单位毫秒（仅在ConfigClient中有效）
		BeatInterval:         5 * 1000,             //心跳间隔时间，单位毫秒（仅在ServiceClient中有效）
		CacheDir:             "./data/nacos/cache", //缓存目录
		LogDir:               "./data/nacos/log",   //日志目录
		UpdateThreadNum:      20,                   //更新服务的线程数
		NotLoadCacheAtStart:  true,                 //在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: true,                 //当服务列表为空时是否更新本地缓存，true--更新,false--不更新
		NamespaceId:          constant.DEFAULT_NAMESPACE_ID,
	}

	serverConfigs := make([]constant.ServerConfig, 0)
	serverConfigs = append(serverConfigs, constant.ServerConfig{
		IpAddr:      "console.nacos.io",
		ContextPath: "/nacos",
		Port:        uint64(80),
	})

	client, err := clients.CreateConfigClient(map[string]interface{}{
		constant.KEY_SERVER_CONFIGS: serverConfigs,
		constant.KEY_CLIENT_CONFIG:  clientConfig,
	})
	if err != nil {
		log.Panic("error create ConfigClient: ", err)
	}

	return client
}

func initIniNacosData(address, group, dataId, tenant string, size, len int) map[string]string {
	m := make(map[string]string)
	content := ""

	for i := 0; i < size; i++ {
		key := "key-" + strconv.Itoa(i)

		for j := 0; j < len; j++ {
			kk := key + "." + "x" + strconv.Itoa(j)
			val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
			k := strings.Replace(kk, "/", ".", -1)
			//fmt.Println(key, k, value)
			m[k] = val
			content += k
			content += "="
			content += val
			content += "\n"
		}
	}
	//url := "http://172.16.1.248:8848/nacos/v1/cs/configs"
	url := "http://console.nacos.io/nacos/v1/cs/configs"
	//buf := strings.NewReader("appName=&namespaceId=&type=properties&dataId=" + dataId + "&group=" + group + "&tenant=" + tenant + "&content=" + content)
	buf := strings.NewReader("type=properties&dataId=" + dataId + "&group=" + group + "&tenant=" + tenant + "&content=" + content)
	//fmt.Println(url, buf)
	res, err := http.Post(url, "application/x-www-form-urlencoded", buf)
	//fmt.Println(res, err)
	data, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(data), err)
	return m

}
