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
		Convey("", func() {
			So(len(y.Keys()), ShouldEqual, 4)
			So(y.Values["key2"], ShouldEqual, "2")
			So(y.Values["application.port"], ShouldEqual, "19002")
		})

	})
}
