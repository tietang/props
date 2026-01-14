package apollo

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"github.com/tietang/props/v3/ini"
	"github.com/tietang/props/v3/kvs"
	"github.com/tietang/props/v3/yam"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultCluster            = "default"
	apolloHeaderAuthorization = "Authorization"
	apolloHeaderTimestamp     = "Timestamp"
	signAuthorizationFormat   = "Apollo %s:%s"
	signDelimiter             = "\n"

	pollInterval          = time.Second * 2
	pollTimeout           = time.Second * 90
	queryTimeout          = time.Second * 5
	defaultNotificationID = -1
)

var _ kvs.ConfigSource = new(ApolloConfigSource)

// 通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type ApolloConfigSource struct {
	kvs.MapProperties
	name                string
	address             string
	appId               string
	namespaces          []string
	cluster             string
	clientIP            string
	secretAccessKey     string
	notifyNamespaces    []*Notification
	configRequester     *requester
	notifyPoolRequester *requester
	contentType         kvs.ContentType
}

func NewApolloConfigSourceWithSecret(address, appId, secretAccessKey string, namespaces []string) *ApolloConfigSource {
	s := newApolloConfigSource(address, appId, namespaces)
	s.secretAccessKey = secretAccessKey
	s.init()
	return s
}

func newApolloConfigSource(address string, appId string, namespaces []string) *ApolloConfigSource {
	s := &ApolloConfigSource{}
	s.name = "apollo:" + address
	s.appId = appId
	s.address = address
	s.namespaces = namespaces
	s.contentType = kvs.KeyValueContentType
	s.clientIP, _ = utils.GetExternalIP()
	s.cluster = "default"
	s.Values = make(map[string]string)
	s.notifyNamespaces = make([]*Notification, 0, 8)
	s.notifyPoolRequester = newRequester(&http.Client{
		Timeout: pollTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	})
	s.configRequester = newRequester(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
		Timeout: queryTimeout,
	})

	return s
}

func NewApolloConfigSource(address, appId string, namespaces []string) *ApolloConfigSource {
	s := newApolloConfigSource(address, appId, namespaces)
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
func NewApolloCompositeConfigSourceWithSecret(url, appId, secretAccessKey string, namespaces []string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "ApolloKevValueWithSecret"
	c := NewApolloConfigSourceWithSecret(url, appId, secretAccessKey, namespaces)
	s.Add(c)
	return s
}

func (s *ApolloConfigSource) init() {
	for _, ns := range s.namespaces {
		s.loadRemoteConfig(ns)
	}
	s.watchALl()
	go s.pollNotifyStart()
}

func (s *ApolloConfigSource) watchALl() {
	s.AddWatchNamespaces(s.namespaces...)
}

func (s *ApolloConfigSource) loadRemoteConfig(ns string) {
	kvs, err := s.GetConfigsFromCache(s.address, s.appId, s.cluster, ns)
	if err != nil {
		log.Error(err)
		return
	}
	for key, value := range kvs {
		s.initValue(ns, key, value)
	}
}
func (p *ApolloConfigSource) AddWatchNamespaces(namespaces ...string) {
	for _, namespace := range namespaces {
		p.notifyNamespaces = append(p.notifyNamespaces, &Notification{
			NamespaceName:  namespace,
			NotificationId: defaultNotificationID,
		})
	}
}
func (p *ApolloConfigSource) RemoveWatchedNamespace(namespace string) {
	for i, ns := range p.notifyNamespaces {
		if ns.NamespaceName == namespace {
			p.notifyNamespaces = append(p.notifyNamespaces[0:i], p.notifyNamespaces[i:]...)
		}
	}
}

func (s *ApolloConfigSource) initValue(namespace, key, value string) {
	contentType := kvs.GetContentTypeByName(namespace)
	if contentType == kvs.ContentUnknown {
		contentType = kvs.ReadContentType(value)
	}
	if contentType == kvs.KeyValueContentType {
		s.Set(key, value)
	} else if contentType == kvs.ContentProps || contentType == kvs.ContentProperties {
		s.Set(key, value)
	} else if contentType == kvs.ContentIni {
		s.findIni(value)
	} else if contentType == kvs.ContentYaml || contentType == kvs.ContentYam || contentType == kvs.ContentYml {
		s.findYaml(value)
	} else {
		s.Set(key, value)
	}
}

func (s *ApolloConfigSource) pollNotifyStart() {
	t2 := time.NewTimer(pollInterval)
	for {
		select {
		case <-t2.C:
			s.poll()
			t2.Reset(pollInterval)
		}
	}
}
func (s *ApolloConfigSource) poll() {
	url := notificationURL(s.address, s.appId, s.cluster, s.clientIP, s.notifyNamespaces)
	data, err := s.notifyPoolRequester.request(url, s.appId, s.secretAccessKey)
	if err != nil {
		log.Error(err)
		return
	}
	notifications := make([]Notification, 0, 8)
	err = json.Unmarshal(data, &notifications)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(string(data))
	for _, notification := range notifications {
		for _, notify := range s.notifyNamespaces {
			if notify.NamespaceName == notification.NamespaceName {
				notify.NotificationId = notification.NotificationId
			}
		}
	}

	for _, notification := range notifications {
		s.loadRemoteConfig(notification.NamespaceName)
	}

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

	return c.configRequester.requestGetKeyValue(url, appID, c.secretAccessKey)

	//req, err := http.NewRequest("GET", url, nil)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//if c.secretAccessKey != "" {
	//	t := getTimestamp()
	//	req.Header.Set("Authorization", getAuthorization(url, t, c.appId, c.secretAccessKey))
	//	req.Header.Set("Timestamp", t)
	//}
	//
	////调用请求
	//res, err := http.DefaultClient.Do(req)
	//
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//// 如果出错就不需要close，因此defer语句放在err处理逻辑后面
	//defer res.Body.Close()
	//if res.StatusCode != http.StatusOK {
	//	msg := "请求错误： status: " + res.Status
	//	log.Error(msg)
	//	return nil, errors.New(msg)
	//}
	////处理response,读取Response body
	//respBody, err := io.ReadAll(res.Body)
	//
	////
	//if err := res.Body.Close(); err != nil {
	//	log.Error(err)
	//}
	//kvs = KeyValue{}
	//err = json.Unmarshal(respBody, &kvs)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//return kvs, err

}

type KeyValue map[string]string

type ConfigRes struct {
	AppID      string   `json:"appId"`          // appId: "AppTest",
	Cluster    string   `json:"cluster"`        // cluster: "default",
	Namespace  string   `json:"namespaceName"`  // namespaceName: "TEST.Namespace1",
	KeyValue   KeyValue `json:"configurations"` // configurations: {Name: "Foo"},
	ReleaseKey string   `json:"releaseKey"`     // releaseKey: "20181017110222-5ce3b2da895720e8"
}

type Notification struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationId int64  `json:"notificationId,omitempty"`
}
