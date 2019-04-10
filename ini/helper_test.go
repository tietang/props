package ini

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestByIni(t *testing.T) {

	Convey("read", t, func() {

		s := `

[key0]

x0.y0=value-00
x0.y1=value-01

[key1]

x1.y0=value-10
x1.y1=value-11

`
		props := ByIni(s)
		So(props, ShouldNotBeNil)
		So(len(props.Keys()), ShouldEqual, 4)

	})
}
