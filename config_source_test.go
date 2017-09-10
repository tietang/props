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
        ps := NewPropertiesConfigSource(file + "/t.test")
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
func TestCompositeConfigSource_Order(t *testing.T) {

    conf := NewEmptyCompositeConfigSource()

    kv1 := []string{"go.app.key1", "value1", "value1-2"}
    kv2 := []string{"go.app.key2", "value2", "value2-2"}

    p1 := NewEmptyMapConfigSource("map1")
    p1.Set(kv1[0], kv1[1])
    p1.Set(kv2[0], kv2[1])
    p2 := NewEmptyMapConfigSource("map2")
    p2.Set(kv1[0], kv1[2])
    p2.Set(kv2[0], kv2[2])
    conf.Add(p1)
    conf.Add(p2)
    envs := os.Environ()
    envLen := len(envs)

    Convey("Test CompositeConfigSource order", t, func() {
        keys := conf.Keys()
        Convey("key len", func() {
            So(len(keys), ShouldEqual, 2+envLen)
        })
        Convey("key/value", func() {
            value1, err := conf.Get(kv1[0])
            So(err, ShouldBeNil)
            So(value1, ShouldEqual, kv1[2])
            value2, err := conf.Get(kv2[0])
            So(err, ShouldBeNil)
            So(value2, ShouldEqual, kv2[2])

        })

    })
}

func TestPlaceholder_String(t *testing.T) {
    p := NewEmptyMapConfigSource("map2")
    p.Set("orign.key1", "v1")
    p.Set("orign.key2", "v2")
    p.Set("ph.key1", "${orign.key1}")
    p.Set("ph.key2", "${orign.key1}:${orign.key2}")
    conf := NewDefaultCompositeConfigSource(p)
    Convey("Test CompositeConfigSource Placeholder", t, func() {
        Convey("ph simple", func() {
            ov1, err := conf.Get("orign.key1")
            So(err, ShouldBeNil)
            phv1, err := conf.Get("ph.key1")
            So(err, ShouldBeNil)
            So(phv1, ShouldEqual, ov1)
        })
        Convey("ph muti composite", func() {
            ov1, err := conf.Get("orign.key1")
            So(err, ShouldBeNil)
            ov2, err := conf.Get("orign.key2")
            So(err, ShouldBeNil)
            phv2, err := conf.Get("ph.key2")
            So(err, ShouldBeNil)
            So(phv2, ShouldEqual, ov1+":"+ov2)
        })
    })

}



func TestPlaceholder_Int(t *testing.T) {
    p := NewEmptyMapConfigSource("map2")
    p.Set("orign.key1", "1")
    p.Set("orign.key2", "2")
    p.Set("ph.key1", "${orign.key1}")
    p.Set("ph.key2", "${orign.key1}+${orign.key2}")
    conf := NewDefaultCompositeConfigSource(p)
    Convey("Test CompositeConfigSource Placeholder", t, func() {
        Convey("ph simple", func() {
            ov1, err := conf.GetInt("orign.key1")
            So(err, ShouldBeNil)
            phv1, err := conf.GetInt("ph.key1")
            So(err, ShouldBeNil)
            So(phv1, ShouldEqual, ov1)
        })
        Convey("ph muti composite", func() {
            ov1, err := conf.Get("orign.key1")
            So(err, ShouldBeNil)
            ov2, err := conf.Get("orign.key2")
            So(err, ShouldBeNil)
            phv2, err := conf.Get("ph.key2")
            So(err, ShouldBeNil)
            So(phv2, ShouldEqual, ov1+"+"+ov2)
        })
    })

}
