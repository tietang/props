package nacos

import (
    "fmt"
    "github.com/hashicorp/consul/api"
    . "github.com/smartystreets/goconvey/convey"
    "net/http"
    "strconv"
    "strings"
    "testing"
)

func TestNacosIniConfigSource(t *testing.T) {
    address := "127.0.0.1:8848"
    size := 10
    inilen := 3

    c := NewNacosPropsConfigSource(address)
    c.DataId = "test.id"
    c.Tenant = "testTenant"
    c.Group = "testGroup"

    m := initIniNacosData(address, c, size, inilen)
    c.init()
    Convey("Nacos kv", t, func() {
        keys := c.Keys()

        So(len(keys), ShouldEqual, size*inilen)
        for _, key := range keys {
            v, ok := m[key]
            //fmt.Println(key)
            v1, err := c.Get(key)
            So(ok, ShouldEqual, true)
            So(err, ShouldBeNil)
            So(v, ShouldEqual, v1)
        }
    })

}

func initIniNacosData(address string, ncs *NacosPropsConfigSource, size, len int) map[string]string {
    config := api.DefaultConfig()
    config.Address = address
    m := make(map[string]string)
    content := ""

    for i := 0; i < size; i++ {
        key := "key-" + strconv.Itoa(i)
        value := ""

        for j := 0; j < len; j++ {
            kk := key + "." + "x" + strconv.Itoa(j)
            val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
            value += "x" + strconv.Itoa(j) + "=" + val + "\n"
            k := strings.Replace(kk, "/", ".", -1)
            //fmt.Println(key, k, value)
            m[k] = val
            content += k
            content += "="
            content += val
            content += "\n"
        }
    }
    url := "http://127.0.0.1:8848/nacos/v1/cs/configs"
    buf := strings.NewReader("dataId=" + ncs.DataId + "&group=" + ncs.Group + "&tenant=" + ncs.Tenant + "&content=" + content)
    fmt.Println(url, buf)
    res,err:=http.Post(url, "application/x-www-form-urlencoded", buf)
    fmt.Println(res,err)
    return m

}
