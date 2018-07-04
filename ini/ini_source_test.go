package ini

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    "github.com/tietang/props/kvs"
)

var inSourceTest *IniFileConfigSourceTest

func init() {
    inSourceTest = &IniFileConfigSourceTest{}
}

type IniFileConfigSourceTest struct {
    IniConfigSourceTest
}

func (c *IniFileConfigSourceTest) CreateConfigSource(p *kvs.Properties) kvs.ConfigSource {
    cs := NewIniFileConfigSource("test-p")
    cs.Values = p.Values
    return cs
}

func TestIniFileConfigSource(t *testing.T) {
    Convey("测试PropertiesConfigSource", t, func() {
        inSourceTest.TestUtils_Properties_GetBool(t)
        inSourceTest.TestUtils_Properties_GetBoolDefault(t)
        inSourceTest.TestUtils_Properties_GetDuration(t)
        inSourceTest.TestUtils_Properties_GetDurationDefault(t)
        inSourceTest.TestUtils_Properties_GetInt(t)
        inSourceTest.TestUtils_Properties_GetIntDefault(t)
        inSourceTest.TestUtils_Read(t)
    })
}
