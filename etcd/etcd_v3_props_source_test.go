// +build go1.9

package etcd

import (
    "testing"
    "strconv"
    . "github.com/smartystreets/goconvey/convey"
    "strings"
    "time"
    "context"
    "github.com/coreos/etcd/clientv3"
    "github.com/coreos/etcd/clientv3/namespace"
    "fmt"
    log "github.com/sirupsen/logrus"
)


func TestEtcdV3PropsConfigSource(t *testing.T) {

    address := testEtcd.Address

    root := "/config101/test/propsdemo1"
    size := 10
    inilen := 3
    m := initEtcdV3PropsData(address, root, size, inilen)
    c := NewEtcdV3PropsConfigSource(address, root)
    Convey("etcd props", t, func() {
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

func initEtcdV3PropsData(address, root string, size, len int) map[string]string {

    cfg := clientv3.Config{
        Endpoints:   []string{address},
        DialTimeout: 3 * time.Second,
    }
    c, err := clientv3.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    //kv := clientv3.NewKV(c)
    kv := namespace.NewKV(c, root)
    r, e := kv.Delete(context.Background(), "/", clientv3.WithPrefix())
    fmt.Println(r, e)
    m := make(map[string]string)
    for i := 0; i < size; i++ {
        key := "/key-" + strconv.Itoa(i)
        //keyFull := path.Join(root, key)
        value := ""

        for j := 0; j < len; j++ {
            kk := key + "." + "x" + strconv.Itoa(j)
            val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
            value += "x" + strconv.Itoa(j) + "=" + val + "\n"
            k := strings.Replace(kk, "/", ".", -1)
            //fmt.Println(key, k, value)
            m[k] = val
        }
        kv.Put(context.Background(), key, value)
    }
    //c.KV = namespace.NewKV(c.KV, root)
    res, err := kv.Get(context.Background(), root, clientv3.WithKeysOnly())
    //res, err := c.KV.Get(context.Background(), "key", clientv3.WithKeysOnly())
    fmt.Println(res.Kvs, res.Count, res.More)
    return m

}
