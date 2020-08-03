package kvs

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDefaultEval_EvalAll(t *testing.T) {
	p := NewMapProperties()
	p.Set("orign.key1", "v1")
	p.Set("orign.key2", "v2")
	p.Set("ph.key1", "${orign.key1}")
	p.Set("ph.key2", "${orign.key1}:${orign.key2}")
	p.Set("ph1.key1", "${ph.key1}")
	p.Set("ph1.key2", "${ph.key2}")
	p.Set("ph2.key1", "${ph1.key1}")
	p.Set("ph2.key2", "${ph1.key2}")

	Convey("TestDefaultEval_EvalAll", t, func() {
		e := NewEval(p)
		e.EvalAll()
		Convey("eval all", func() {
			So(p.GetDefault("ph.key1", ""), ShouldEqual, p.GetDefault("ph.key1", ""))
			So(p.GetDefault("ph.key2", ""), ShouldEqual, p.GetDefault("ph.key2", ""))
			So(p.GetDefault("ph1.key1", ""), ShouldEqual, p.GetDefault("ph.key1", ""))
			So(p.GetDefault("ph1.key2", ""), ShouldEqual, p.GetDefault("ph.key2", ""))
			So(p.GetDefault("ph2.key1", ""), ShouldEqual, p.GetDefault("ph.key1", ""))
			So(p.GetDefault("ph2.key2", ""), ShouldEqual, p.GetDefault("ph.key2", ""))
			//
			So(p.GetDefault("ph.key1", ""), ShouldEqual, p.GetDefault("ph1.key1", ""))
			So(p.GetDefault("ph.key2", ""), ShouldEqual, p.GetDefault("ph1.key2", ""))
			So(p.GetDefault("ph2.key1", ""), ShouldEqual, p.GetDefault("ph1.key1", ""))
			So(p.GetDefault("ph2.key2", ""), ShouldEqual, p.GetDefault("ph1.key2", ""))
		})
	})

}
