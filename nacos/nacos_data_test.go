package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func update(address, group, dataId, namespaceId string, size, len int) map[string]string {
	m := make(map[string]string)
	c := client(address, namespaceId)
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
		Content: content,
	}
	b, e := c.PublishConfig(vo)
	log.Info("updated: ", b, e)

	return m
}

func client(address, namespaceId string) config_client.IConfigClient {
	clientConfig := constant.ClientConfig{
		TimeoutMs:            10 * 1000,       //请求Nacos服务端的超时时间，默认是10000ms
		BeatInterval:         5 * 1000,        //心跳间隔时间，单位毫秒（仅在ServiceClient中有效）
		CacheDir:             "./nacos/cache", //缓存目录
		LogDir:               "./nacos/log",   //日志目录
		UpdateThreadNum:      20,              //更新服务的线程数
		NotLoadCacheAtStart:  true,            //在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: false,           //当服务列表为空时是否更新本地缓存，true--更新,false--不更新,当service返回的实例列表为空时，不更新缓存，用于推空保护
		RotateTime:           "1h",            // 日志轮转周期，比如：30m, 1h, 24h, 默认是24h
		MaxAge:               3,               // 日志最大文件数，默认3
		NamespaceId:          namespaceId,
	}
	a := strings.Split(address, ":")
	port := 80
	if len(a) == 2 {
		var err error
		port, err = strconv.Atoi(a[1])
		if err != nil {
			log.Error("error config nacos address:", address)
		}
	}

	serverConfigs := make([]constant.ServerConfig, 0)
	serverConfigs = append(serverConfigs, constant.ServerConfig{
		IpAddr:      a[0],
		ContextPath: "/nacos",
		Port:        uint64(port),
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

func initIniNacosData(address, group, dataId, namespaceId string, size, length int) map[string]string {
	m := make(map[string]string)
	content := ""
	c := client(address, namespaceId)
	for i := 0; i < size; i++ {
		key := "key-" + strconv.Itoa(i)

		for j := 0; j < length; j++ {
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

	//url := fmt.Sprintf("http://%s/nacos/v1/cs/configs", address)
	////buf := strings.NewReader("appName=&namespaceId=&type=properties&dataId=" + dataId + "&group=" + group + "&tenant=" + tenant + "&content=" + content)
	//buf := strings.NewReader("type=properties&dataId=" + dataId + "&group=" + group + "&tenant=" + tenant + "&content=" + content)
	////fmt.Println(url, buf)
	//res, err := http.Post(url, "application/x-www-form-urlencoded", buf)
	////fmt.Println(res, err)
	//data, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(data), err)
	vo := vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	}
	ok, err := c.PublishConfig(vo)
	//log.Info(content)
	log.Info("add: ", ok, err)
	log.Info(len(m))

	return m

}
