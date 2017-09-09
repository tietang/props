package props

import (
    "testing"
    "github.com/hashicorp/consul/api"
    "strconv"
    . "github.com/smartystreets/goconvey/convey"
    "path"
    "strings"
)

func TestConsulIniConfigSource(t *testing.T) {
    address := "127.0.0.1:8500"
    //address := "172.16.1.248:8500"
    root := "config101/test/inidemo"
    size := 10
    inilen := 3
    m := initIniConsulData(address, root, size, inilen)
    c := NewConsulIniConfigSource(address, root)
    Convey("consul kv", t, func() {
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

func initIniConsulData(address, root string, size, len int) map[string]string {
    config := api.DefaultConfig()
    config.Address = address
    client, err := api.NewClient(config)
    if err != nil {
        panic(err)
    }
    m := make(map[string]string)
    kv := client.KV()
    wq := &api.WriteOptions{}

    for i := 0; i < size; i++ {
        key := "key-" + strconv.Itoa(i)
        keyFull := path.Join(root, key)

        value := ""

        for j := 0; j < len; j++ {
            kk := key + "." + "x" + strconv.Itoa(j)
            val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
            value += "x" + strconv.Itoa(j) + "=" + val + "\n"
            k := strings.Replace(kk, "/", ".", -1)
            //fmt.Println(key, k, value)
            m[k] = val
        }
        kvp := &api.KVPair{
            Key:   keyFull,
            Value: []byte(value),
        }
        kv.Put(kvp, wq)

    }

    return m

}
