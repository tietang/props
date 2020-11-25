package kvs

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const (
	STR_VAL          = "str demo"
	INT_VAL          = 122
	INT_VAL_STR      = "122"
	DURATION_VAL     = time.Duration(2) * time.Second
	DURATION_VAL_STR = "2s"
	BOOL_VAL         = true
	BOOL_VAL_STR     = "true"
)

func TestStruct(t *testing.T) {
	type PlatStruct struct {
		StrVal      string
		IntVal      int
		DurationVal time.Duration
		BoolVal     bool
	}
	ps := NewMapProperties()
	ps.Set("ums.strVal", STR_VAL)
	ps.Set("ums.intVal", INT_VAL_STR)
	ps.Set("ums.durationVal", DURATION_VAL_STR)
	ps.Set("ums.boolVal", BOOL_VAL_STR)

	Convey("TestStructUnmarshal", t, func() {

		Convey("test flat struct unmarshal", func() {
			s := &PlatStruct{}
			err := Unmarshal(ps, s, "ums")
			So(err, ShouldBeNil)
			So(s.StrVal, ShouldEqual, STR_VAL)
			So(s.IntVal, ShouldEqual, INT_VAL)
			So(s.DurationVal, ShouldEqual, DURATION_VAL)
			So(s.BoolVal, ShouldEqual, BOOL_VAL)
		})

	})
}

func TestInnerStruct(t *testing.T) {
	type PlatStruct struct {
		StrVal      string
		IntVal      int
		DurationVal time.Duration
		BoolVal     bool
	}
	type OuterStruct struct {
		Inner PlatStruct
	}
	ps := NewMapProperties()
	ps.Set("ums.inner.strVal", STR_VAL)
	ps.Set("ums.inner.intVal", INT_VAL_STR)
	ps.Set("ums.inner.durationVal", DURATION_VAL_STR)
	ps.Set("ums.inner.boolVal", BOOL_VAL_STR)

	Convey("TestStructUnmarshal", t, func() {

		Convey("test inner struct unmarshal", func() {
			s := &OuterStruct{}
			err := Unmarshal(ps, s, "ums")
			So(err, ShouldBeNil)
			So(s.Inner.StrVal, ShouldEqual, STR_VAL)
			So(s.Inner.IntVal, ShouldEqual, INT_VAL)
			So(s.Inner.DurationVal, ShouldEqual, DURATION_VAL)
			So(s.Inner.BoolVal, ShouldEqual, BOOL_VAL)
		})

	})
}

func TestAnonymityStruct(t *testing.T) {
	type PlatStruct struct {
		StrVal      string
		IntVal      int
		DurationVal time.Duration
		BoolVal     bool
	}
	type OuterStruct struct {
		PlatStruct
	}
	ps := NewMapProperties()
	ps.Set("ums.strVal", STR_VAL)
	ps.Set("ums.intVal", INT_VAL_STR)
	ps.Set("ums.durationVal", DURATION_VAL_STR)
	ps.Set("ums.boolVal", BOOL_VAL_STR)

	Convey("TestStructUnmarshal", t, func() {

		Convey("test anonymity struct unmarshal", func() {
			s := &OuterStruct{}
			err := Unmarshal(ps, s, "ums")
			So(err, ShouldBeNil)
			So(s.StrVal, ShouldEqual, STR_VAL)
			So(s.IntVal, ShouldEqual, INT_VAL)
			So(s.DurationVal, ShouldEqual, DURATION_VAL)
			So(s.BoolVal, ShouldEqual, BOOL_VAL)
		})

	})
}

func TestNestStruct(t *testing.T) {

	type OuterStruct struct {
		Inner struct {
			StrVal      string
			IntVal      int
			DurationVal time.Duration
			BoolVal     bool
		}
	}
	ps := NewMapProperties()
	ps.Set("ums.inner.strVal", STR_VAL)
	ps.Set("ums.inner.intVal", INT_VAL_STR)
	ps.Set("ums.inner.durationVal", DURATION_VAL_STR)
	ps.Set("ums.inner.boolVal", BOOL_VAL_STR)

	Convey("TestStructUnmarshal", t, func() {

		Convey("test nest struct unmarshal", func() {
			s := &OuterStruct{}
			err := Unmarshal(ps, s, "ums")
			So(err, ShouldBeNil)
			So(s.Inner.StrVal, ShouldEqual, STR_VAL)
			So(s.Inner.IntVal, ShouldEqual, INT_VAL)
			So(s.Inner.DurationVal, ShouldEqual, DURATION_VAL)
			So(s.Inner.BoolVal, ShouldEqual, BOOL_VAL)
		})

	})
}

func TestAnonymityNestStruct(t *testing.T) {
	type PlatStruct struct {
		Inner struct {
			StrVal      string
			IntVal      int
			DurationVal time.Duration
			BoolVal     bool
		}
	}

	type OuterStruct struct {
		PlatStruct
	}
	ps := NewMapProperties()
	ps.Set("ums.inner.strVal", STR_VAL)
	ps.Set("ums.inner.intVal", INT_VAL_STR)
	ps.Set("ums.inner.durationVal", DURATION_VAL_STR)
	ps.Set("ums.inner.boolVal", BOOL_VAL_STR)

	Convey("TestStructUnmarshal", t, func() {

		Convey("test Anonymity and Nest struct unmarshal", func() {
			s := &OuterStruct{}
			err := Unmarshal(ps, s, "ums")
			So(err, ShouldBeNil)
			So(s.Inner.StrVal, ShouldEqual, STR_VAL)
			So(s.Inner.IntVal, ShouldEqual, INT_VAL)
			So(s.Inner.DurationVal, ShouldEqual, DURATION_VAL)
			So(s.Inner.BoolVal, ShouldEqual, BOOL_VAL)
		})

	})
}
