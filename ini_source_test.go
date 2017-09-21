package props

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

var inSourceTest *IniFileConfigSourceTest

func init() {
    inSourceTest = &IniFileConfigSourceTest{}
}

type IniFileConfigSourceTest struct {
    ConfigSourceTest
}

func (c *IniFileConfigSourceTest) CreateConfigSource(p *Properties) ConfigSource {
    cs := NewIniFileConfigSource("test-p")
    cs.values = p.values
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
