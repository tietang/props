package props

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    "os"
)

func TestNewPropertiesConfigSourceByFile(t *testing.T) {
    kv1 := []string{"go.app.key1", "value1"}
    kv2 := []string{"go.app.key2", "value2"}
    Convey("TestNewPropertiesConfigSourceByFile", t, func() {
        file, _ := os.Getwd()
        ps := NewPropertiesConfigSourceByFile("props", file+"/t.test")
        keys := ps.Keys()
        Convey("key len", func() {
            So(len(keys), ShouldEqual, 2+1)
        })
        Convey("key/value", func() {
            for k, v := range ps.values {
                if k == kv1[0] {
                    So(kv1[1], ShouldEqual, v)
                }
                if k == kv2[0] {
                    So(kv2[1], ShouldEqual, v)
                }

            }
        })

    })
}
