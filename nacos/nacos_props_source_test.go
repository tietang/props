package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
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
	//address := "172.16.1.248:8848"
	//http://console.nacos.io/nacos/v1/cs/configs?show=all&dataId=q123&group=DEFAULT_GROUP&tenant=&namespaceId=
	//http://console.nacos.io/nacos/v1/cs/configs
	address := "console.nacos.io"

	size := 10
	inilen := 3
	dataId := "test.id"
	tenant := constant.DEFAULT_NAMESPACE_ID
	group := "testGroup"
	m := initIniNacosData(address, group, dataId, tenant, size, inilen)
	c := NewNacosPropsConfigSource(address, group, dataId, tenant)

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
		m = update(address, group, dataId, tenant, size, inilen)
		time.Sleep(time.Second * 5)
		keys = c.Keys()
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

	time.Sleep(time.Second * 20)

}
