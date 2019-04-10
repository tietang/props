package yam

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

var s = `
application:
  name: go-example
  port: 19002
key1:
key2: 2
`

func TestNewYamlProperties(t *testing.T) {
	y := NewYamlProperties()

	Convey("测试NewYamlProperties", t, func() {
		err := y.Load(strings.NewReader(s))
		So(err, ShouldBeNil)
		Convey("正常存在的key", func() {
			So(len(y.Keys()), ShouldEqual, 4)

			v, ok := y.Values["key1"]
			So(ok, ShouldBeTrue)
			So(v, ShouldBeEmpty)
			So(y.Values["key2"], ShouldEqual, "2")
			So(y.Values["application.port"], ShouldEqual, "19002")

			//
			v, err := y.Get("key1")
			So(err, ShouldBeNil)
			So(v, ShouldBeEmpty)
			v, err = y.Get("key2")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "2")

			v, err = y.Get("application.port")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "19002")

			i, err := y.GetInt("application.port")
			So(err, ShouldBeNil)
			So(i, ShouldEqual, 19002)

		})
		Convey("不存在的key", func() {
			So(len(y.Keys()), ShouldEqual, 4)

			v, ok := y.Values["key1-3213131242nfdksfdks"]
			So(ok, ShouldBeFalse)
			So(v, ShouldBeEmpty)

			v, err := y.Get("key1-3213131242nfdksfdks")
			So(err, ShouldNotBeNil)
			So(v, ShouldBeEmpty)

		})

	})
}
