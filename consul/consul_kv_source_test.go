package consul

import (
    "testing"
    "github.com/hashicorp/consul/api"
    "strconv"
    . "github.com/smartystreets/goconvey/convey"
    "path"
    "strings"
)



func TestConsulKeyValueConfigSource(t *testing.T) {

    address := testConsul.Address

    root := "config101/test/kvdemo1"
    size := 10
    m := initConsulData(address, root, size)
    c := NewConsulKeyValueConfigSource(address, root)
    Convey("consul kv", t, func() {
        keys := c.Keys()
        So(len(keys), ShouldEqual, size)
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

func initConsulData(address, root string, size int) map[string]string {
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
        key := "key/x" + strconv.Itoa(i)
        keyFull := path.Join(root, key)
        value := "value-" + strconv.Itoa(i)
        kvp := &api.KVPair{
            Key:   keyFull,
            Value: []byte(value),
        }
        kv.Put(kvp, wq)
        k := strings.Replace(key, "/", ".", -1)
        //fmt.Println(key, k, value)
        m[k] = value
    }

    return m

}
