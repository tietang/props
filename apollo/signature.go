package apollo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	neturl "net/url"
	"time"
)

const DELIMITER = "\n"

func getAuthorization(url, timestamp, appid, accesskeySecret string) string {
	sign := timestamp + DELIMITER + getURL2PathWithQuery(url)
	key := []byte(accesskeySecret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(sign))
	ss := mac.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(ss)
	return fmt.Sprintf(
		signAuthorizationFormat,
		appid,
		signature,
	)
	//return fmt.Sprintf("Apollo %s:%s", appid, signature)
}

func getTimestamp() string {
	t := time.Now().UnixMilli() //time.Now().UnixNano() / 1e6
	return fmt.Sprintf("%d", t)
}

func getURL2PathWithQuery(rawurl string) string {
	url, _ := neturl.Parse(rawurl)
	return url.RequestURI()
}
