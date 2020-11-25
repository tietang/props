package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

//http://console.nacos.io/nacos/v1/cs/configs
//var address = "console.nacos.io"
var address = "172.16.1.248:8848"

func TestNacosClientIniConfigSource2(t *testing.T) {

	//http://console.nacos.io/nacos/v1/cs/configs?
	// show=all&dataId=xxx123&group=DEFAULT_GROUP&tenant=&namespaceId=
	dataId := "xxx123"
	namespaceId := constant.DEFAULT_NAMESPACE_ID
	group := constant.DEFAULT_GROUP
	c := NewNacosClientPropsConfigSource(address, group, dataId, namespaceId)
	fmt.Println(c.Keys())
}

func TestNacosClientIniConfigSource(t *testing.T) {
	//http://console.nacos.io/nacos/v1/cs/configs?show=all&dataId=q123&group=DEFAULT_GROUP&tenant=&namespaceId=

	size := 10
	inilen := 3
	dataId := "test1"
	namespaceId := "" //constant.DEFAULT_NAMESPACE_ID
	group := "testGroup"
	//group := constant.DEFAULT_GROUP
	m := initIniNacosData(address, group, dataId, namespaceId, size, inilen)
	c := NewNacosClientPropsConfigSource(address, group, dataId, namespaceId)

	Convey("Nacos kv", t, func() {
		keys := c.Keys()
		//fmt.Println(keys)

		So(len(keys), ShouldEqual, size*inilen)
		for _, key := range keys {
			v, ok := m[key]
			//fmt.Println(key)
			v1, err := c.Get(key)
			So(ok, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, v1)
		}
		m = update(address, group, dataId, namespaceId, size, inilen)
		time.Sleep(time.Second * 10)
		keys = c.Keys()
		So(len(keys), ShouldEqual, size*inilen)
		for _, key := range keys {
			v, ok := m[key]
			//fmt.Println(key)
			v1, err := c.Get(key)
			So(ok, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(v1, ShouldEqual, v)
		}
	})

}
