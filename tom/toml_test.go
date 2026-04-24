package tom

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var s = `
key1 = ""
key2 = "d"

[application]
name = "go-jiebaexample"
port =  19002

`

func TestNewTomlProperties(t *testing.T) {
	y := NewTomlProperties()

	Convey("测试NewTomlProperties", t, func() {
		err := y.Load(strings.NewReader(s))
		So(err, ShouldBeNil)
		Convey("正常存在的key", func() {
			So(len(y.Keys()), ShouldEqual, 4)

			v, ok := y.Values["key1"]
			So(ok, ShouldBeTrue)
			So(v, ShouldBeEmpty)
			So(y.Values["key2"], ShouldEqual, "d")
			So(y.Values["application.port"], ShouldEqual, "19002")

			//
			v, err := y.Get("key1")
			So(err, ShouldBeNil)
			So(v, ShouldBeEmpty)
			v, err = y.Get("key2")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "d")

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
