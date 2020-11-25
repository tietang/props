package kvs

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestMapStruct(t *testing.T) {
	type PlatStruct struct {
		StrVal      string
		IntVal      int
		DurationVal time.Duration
		BoolVal     bool
	}
	ps := NewMapProperties()
	ps.Set("ums.test1.strVal", STR_VAL)
	ps.Set("ums.test1.intVal", INT_VAL_STR)
	ps.Set("ums.test1.durationVal", DURATION_VAL_STR)
	ps.Set("ums.test1.boolVal", BOOL_VAL_STR)

	ps.Set("ums.test2.strVal", STR_VAL)
	ps.Set("ums.test2.intVal", INT_VAL_STR)
	ps.Set("ums.test2.durationVal", DURATION_VAL_STR)
	ps.Set("ums.test2.boolVal", BOOL_VAL_STR)

	Convey("TestStructUnmarshal", t, func() {

		Convey("test map struct  Unmarshal", func() {
			m := make(map[string]*PlatStruct, 0)
			err := Unmarshal(ps, m, "ums")
			So(err, ShouldBeNil)
			s1, ok := m["test1"]
			So(ok, ShouldBeTrue)
			So(s1, ShouldNotBeNil)
			So(s1.StrVal, ShouldEqual, STR_VAL)
			So(s1.IntVal, ShouldEqual, INT_VAL)
			So(s1.DurationVal, ShouldEqual, DURATION_VAL)
			So(s1.BoolVal, ShouldEqual, BOOL_VAL)

			s2, ok := m["test2"]
			So(ok, ShouldBeTrue)
			So(s2, ShouldNotBeNil)
			So(s2.StrVal, ShouldEqual, STR_VAL)
			So(s2.IntVal, ShouldEqual, INT_VAL)
			So(s2.DurationVal, ShouldEqual, DURATION_VAL)
			So(s2.BoolVal, ShouldEqual, BOOL_VAL)

		})

	})
}
