package kvs

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var csTest *PropsConfigSourceTest

func init() {
	csTest = &PropsConfigSourceTest{}
}

type PropsConfigSourceTest struct {
	ConfigSourceTest
}

func (c *PropsConfigSourceTest) CreateConfigSource(p *Properties) ConfigSource {
	cs := NewPropertiesConfigSourceByMap("test-p", p.Values)
	return cs
}

func TestPropertiesConfigSource(t *testing.T) {
	Convey("测试PropertiesConfigSource", t, func() {
		csTest.TestUtils_Properties_GetBool(t)
		csTest.TestUtils_Properties_GetBoolDefault(t)
		csTest.TestUtils_Properties_GetDuration(t)
		csTest.TestUtils_Properties_GetDurationDefault(t)
		csTest.TestUtils_Properties_GetInt(t)
		csTest.TestUtils_Properties_GetIntDefault(t)
		csTest.TestUtils_Read(t)
	})
}
