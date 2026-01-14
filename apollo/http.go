package apollo

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type requester struct {
	client *http.Client
}

func newRequester(client *http.Client) *requester {
	return &requester{client: client}
}

func (r *requester) requestGetKeyValue(url, appid, secretAccessKey string) (kvs KeyValue, err error) {
	respBody, err := r.request(url, appid, secretAccessKey)
	if err != nil {
		return nil, err
	}
	kvs = KeyValue{}
	err = json.Unmarshal(respBody, &kvs)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return kvs, err
}

func (r *requester) request(url string, appid string, secretAccessKey string) (data []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if secretAccessKey != "" {
		timestamp := getTimestamp()
		req.Header.Set(apolloHeaderAuthorization,
			getAuthorization(url, timestamp, appid, secretAccessKey),
		)
		req.Header.Set(apolloHeaderTimestamp, timestamp)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	// 如果出错就不需要close，因此defer语句放在err处理逻辑后面
	defer resp.Body.Close()
	//处理response,读取Response body
	respBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		msg := "请求错误： status: " + resp.Status
		log.Error(msg)
		return nil, errors.New(msg)
	}

	return respBody, nil
}

func httpPrefix(address string) string {
	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		return address
	}

	return fmt.Sprintf("http://%s", address)
}

func notificationURL(address, appID, cluster, clientIP string, notifications []*Notification) string {
	notificationsJson, err := json.Marshal(&notifications)
	if err != nil {
		return ""
	}
	notificationsJsonStr := string(notificationsJson)
	return fmt.Sprintf("%s/notifications/v2?appId=%s&cluster=%s&notifications=%s&ip=%s",
		httpPrefix(address),
		url.QueryEscape(appID),
		url.QueryEscape(cluster),
		url.QueryEscape(notificationsJsonStr),
		clientIP)
}

func configURLWithCache(address, appID, cluster, namespace, clientIP string) string {
	url := fmt.Sprintf("%s/configfiles/json/%s/%s/%s?ip=%s",
		httpPrefix(address),
		url.QueryEscape(appID),
		url.QueryEscape(cluster),
		url.QueryEscape(namespace),
		clientIP,
	)
	return url
}

func configURL(address, appID, cluster, namespace, releaseKey, clientIP string) string {
	url := fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		httpPrefix(address),
		url.QueryEscape(appID),
		url.QueryEscape(cluster),
		url.QueryEscape(namespace),
		releaseKey,
		clientIP,
	)
	return url
}
