package yam

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestNewYamlConfigSource(t *testing.T) {
	y := NewYamlConfigSourceByReader("string", strings.NewReader(s))
	Convey("测试NewYamlConfigSource", t, func() {
		Convey("", func() {
			So(len(y.Keys()), ShouldEqual, 4)
			So(y.Values["key2"], ShouldEqual, "2")
			So(y.Values["application.port"], ShouldEqual, "19002")
		})

	})
}
