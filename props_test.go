package props

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    "strings"
)

func TestRead(t *testing.T) {

    Convey("properties 文件载入", t, func() {
        r := strings.NewReader(`
        k1=v1
        k2:v2
        k3=v3:v3-1
        k4:k4-1=v4
        k5=v5\rk6=v6\nk7=v7\fk8=v9
        #注释

              #空格注释
            \n\t\r


         fjdskl
         ewjkwl
          k10 =  v10  v10-1
         k11=v11-1
         k11=v11-2
        `)
        p, err := readProperties(r)
        //fmt.Println(p.Keys())
        So(err, ShouldBeNil)
        So(p, ShouldNotBeNil)
        Convey("验证=", func() {
            v ,_:= p.Get("k1")
            So(v, ShouldEqual, "v1")
        })
        Convey("验证:", func() {
            v ,_:= p.Get("k2")
            So(v, ShouldEqual, "v2")

        })
        Convey("验证=:优先级1", func() {
            v ,_:= p.Get("k3")
            So(v, ShouldEqual, "v3:v3-1")
        })
        Convey("验证=:优先级2-异常", func() {
            v ,_:= p.Get("k4")
            So(v, ShouldNotBeNil)
        })
        Convey("验证=:优先级2-正常", func() {
            v ,_:= p.Get("k4:k4-1")
            So(v, ShouldEqual, "v4")
        })
        Convey("验证\\r-NOT-EQ", func() {
            v ,_:= p.Get("k5")
            So(v, ShouldNotEqual, "v5")
        })
        Convey("验证\\r-EQ", func() {
            v ,_:= p.Get("k5")
            So(v, ShouldEqual, "v5\\rk6=v6\\nk7=v7\\fk8=v9")
        })
        Convey("验证\\n", func() {
            v ,_:= p.Get("k6")
            So(v, ShouldNotEqual, "v6")
            So(v, ShouldEqual, "")
        })
        Convey("验证\\f", func() {
            v ,_:= p.Get("k7")
            So(v, ShouldNotEqual, "v7")
            So(v, ShouldEqual, "")
        })
        Convey("验证\\r\\n\\f", func() {
            v ,_:= p.Get("k8")
            So(v, ShouldNotEqual, "v8")
            So(v, ShouldEqual, "")

        })
        Convey("验证trimspace", func() {
            v ,_:= p.Get("k10")
            So(v, ShouldEqual, "v10  v10-1")

        })
        Convey("验证key覆盖", func() {
            v ,_:= p.Get("k11")
            So(v, ShouldEqual, "v11-2")

        })

    })

}

func TestProperties_GetBool(t *testing.T) {

    Convey("测试get bool", t, func() {
        r := strings.NewReader(`
        k1=true
        k2:false
        k3=t
        k4=T
        k5=1
        k6=f
        k7=F
        k8=0
        k9=t0
        k10=12
        k11=-12
        `)
        p, err := readProperties(r)
        So(err, ShouldBeNil)
        So(p, ShouldNotBeNil)
        Convey("k1(true) is true", func() {
            v ,_:= p.GetBool("k1")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeTrue)
        })
        Convey("k2(false) is false", func() {
            v ,_:= p.GetBool("k2")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeFalse)
        })
        Convey("k3(t) is true", func() {
            v ,_:= p.GetBool("k3")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeTrue)
        })
        Convey("k4(T) is true", func() {
            v ,_:= p.GetBool("k4")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeTrue)
        })
        Convey("k5(1) is true", func() {
            v ,_:= p.GetBool("k5")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeTrue)
        })

        Convey("k6(f) is false", func() {
            v ,_:= p.GetBool("k6")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeFalse)
        })
        Convey("k7(F) is false", func() {
            v ,_:= p.GetBool("k7")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeFalse)
        })
        Convey("k8(0) is false", func() {
            v ,_:= p.GetBool("k8")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeFalse)
        })
        Convey("k9(t0) is not bool", func() {
            v ,_:= p.GetBool("k9")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeFalse)
        })
        Convey("k10(12) is not bool", func() {
            v ,_:= p.GetBool("k10")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeFalse)
        })
        Convey("k11(-12) is not bool", func() {
            v ,_:= p.GetBool("k10")
            So(v, ShouldNotBeNil)
            So(v, ShouldBeFalse)
        })

    })
}
func TestProperties_GetBoolDefalut(t *testing.T) {
    defaultValue := true
    Convey("测试get bool", t, func() {
        r := strings.NewReader(`
        k9=t0
        k10=12
        k11=-12
        `)
        p, err := readProperties(r)
        So(err, ShouldBeNil)
        So(p, ShouldNotBeNil)

        Convey("k9(t0) is not bool", func() {
            v ,_:= p.GetBoolDefault("k9", defaultValue)
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, defaultValue)
        })
        Convey("k10(12) is not bool", func() {
            v ,_:= p.GetBoolDefault("k10", defaultValue)
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, defaultValue)
        })
        Convey("k11(-12) is not bool", func() {
            v ,_:= p.GetBoolDefault("k10", defaultValue)
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, defaultValue)
        })
        Convey("k12 is not exits", func() {
            v ,_:= p.GetBoolDefault("k12", defaultValue)
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, defaultValue)
        })

    })
}

func TestProperties_GetInt(t *testing.T) {
    Convey("测试get bool", t, func() {
        r := strings.NewReader(`
        k1= 1
        k2: 2
        k3= -1
        k9= t0
        `)
        p, err := readProperties(r)
        So(err, ShouldBeNil)
        So(p, ShouldNotBeNil)
        Convey("k1(1) is 1", func() {
            v ,_:= p.GetInt("k1")
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, 1)
        })
        Convey("k2(2) is 2", func() {
            v ,_:= p.GetInt("k2")
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, 2)
        })
        Convey("k3(-1) is -1", func() {
            v ,_:= p.GetInt("k3")
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, -1)
        })

        Convey("k9(t0) is not int, is 0", func() {
            v ,_:= p.GetInt("k9")
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, 0)
        })
        Convey("k9-NOT(t0) is not exits, is 0", func() {
            v ,_:= p.GetInt("k9-NOT")
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, 0)
        })

    })
}

func TestProperties_GetIntDefault(t *testing.T) {
    defaultValue := -1000
    Convey("测试get bool", t, func() {
        r := strings.NewReader(`
        k9=t0
        k10=-01w
        `)
        p, err := readProperties(r)
        So(err, ShouldBeNil)
        So(p, ShouldNotBeNil)

        Convey("k9(t0) is not int", func() {
            v ,_:= p.GetIntDefault("k9", defaultValue)
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, defaultValue)
        })
        Convey("k10(-01w) is not int", func() {
            v ,_:= p.GetIntDefault("k10", defaultValue)
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, defaultValue)
        })

        Convey("k12 is not exits", func() {
            v ,_:= p.GetIntDefault("k12", defaultValue)
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, defaultValue)
        })
        Convey("k13 is not exits", func() {
            defaultValue = 3222
            v ,_:= p.GetIntDefault("k13", defaultValue)
            So(v, ShouldNotBeNil)
            So(v, ShouldEqual, defaultValue)
        })

    })
}
