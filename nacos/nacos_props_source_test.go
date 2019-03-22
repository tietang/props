package nacos

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestNacosIniConfigSource2(t *testing.T) {
	address := "console.nacos.io"
	//http://console.nacos.io/nacos/v1/cs/configs?
	// show=all&dataId=xxx123&group=DEFAULT_GROUP&tenant=&namespaceId=
	dataId := "xxx123"
	tenant := ""
	group := "DEFAULT_GROUP"
	c := NewNacosPropsConfigSource(address, group, dataId, tenant)
	fmt.Println(c.Keys())
}
func TestNacosIniConfigSource(t *testing.T) {
	address := "172.16.1.248:8848"
	//http://console.nacos.io/nacos/v1/cs/configs?show=all&dataId=q123&group=DEFAULT_GROUP&tenant=&namespaceId=
	//http://console.nacos.io/nacos/v1/cs/configs
	//address := "console.nacos.io"

	size := 10
	inilen := 3
	dataId := "test.id"
	tenant := "testTenant"
	group := "testGroup"
	m := initIniNacosData(address, group, dataId, tenant, size, inilen)
	c := NewNacosPropsConfigSource(address, group, dataId, tenant)

	Convey("Nacos kv", t, func() {
		keys := c.Keys()
		fmt.Println(keys)

		So(len(keys), ShouldEqual, size*inilen)
		for _, key := range keys {
			v, ok := m[key]
			//fmt.Println(key)
			v1, err := c.Get(key)
			So(ok, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, v1)
		}
	})

}

func initIniNacosData(address, group, dataId, tenant string, size, len int) map[string]string {
	m := make(map[string]string)
	content := ""

	for i := 0; i < size; i++ {
		key := "key-" + strconv.Itoa(i)
		value := ""

		for j := 0; j < len; j++ {
			kk := key + "." + "x" + strconv.Itoa(j)
			val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
			value += "x" + strconv.Itoa(j) + "=" + val + "\n"
			k := strings.Replace(kk, "/", ".", -1)
			//fmt.Println(key, k, value)
			m[k] = val
			content += k
			content += "="
			content += val
			content += "\n"
		}
	}
	url := "http://172.16.1.248:8848/nacos/v1/cs/configs"
	//url := "http://console.nacos.io/nacos/v1/cs/configs"
	buf := strings.NewReader("appName=&namespaceId=&type=properties&dataId=" + dataId + "&group=" + group + "&tenant=" + tenant + "&content=" + content)
	fmt.Println(url, buf)
	res, err := http.Post(url, "application/x-www-form-urlencoded", buf)
	fmt.Println(res, err)
	data, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(data), err)
	return m

}
