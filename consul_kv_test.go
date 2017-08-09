package props

import (
    "testing"
    "github.com/hashicorp/consul/api"
    "strconv"
    . "github.com/smartystreets/goconvey/convey"
    "path"
)

func TestConsulKeyValueConfigSource(t *testing.T) {
    //address := "127.0.0.1:8500"
    address := "172.16.1.248:8500"
    root := "config101/test/demo1"
    size := 10
    initConsulData(address, root, size)
    c := NewConsulKeyValueConfigSource("consul", address, root)
    Convey("consul kv", t, func() {
        keys := c.Keys()
        So(len(keys), ShouldEqual, size)
    })

}

func initConsulData(address, root string, size int) {
    config := api.DefaultConfig()
    config.Address = address
    client, err := api.NewClient(config)
    if err != nil {
        panic(err)
    }
    kv := client.KV()
    wq := &api.WriteOptions{}
    for i := 0; i < size; i++ {
        kvp := &api.KVPair{
            Key:   path.Join(root, "key/x"+strconv.Itoa(i)),
            Value: []byte("value-" + strconv.Itoa(i)),
        }
        kv.Put(kvp, wq)
    }

}
