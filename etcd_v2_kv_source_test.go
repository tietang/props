package props

import (
    "testing"
    "strconv"
    . "github.com/smartystreets/goconvey/convey"
    "path"
    "strings"
    "time"
    "github.com/coreos/etcd/client"
    "log"
    "context"
)

var etcd_mock_started = false
var etcdAddress string

func init() {
    //etcdAddress = "http://172.16.1.248:2379"
    etcdAddress = "http://127.0.0.1:2379"
    GetOrNewMockTestEtcd(etcdAddress)
    if !etcd_mock_started {
        go testEtcd.StartMockEtcd()
    }
    testEtcd.WaitingForEtcdStarted()
}

func TestEtcdKeyValueConfigSource(t *testing.T) {

    address := testEtcd.Address

    root := "/config101/test/kvdemo1"
    size := 10
    m := initEtcdData(address, root, size)
    c := NewEtcdKeyValueConfigSource(address, root)
    Convey("etcd kv", t, func() {
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

func initEtcdData(address, root string, size int) map[string]string {

    cfg := client.Config{
        Endpoints: []string{address},
        Transport: client.DefaultTransport,
        // set timeout per request to fail fast when the target endpoint is unavailable
        HeaderTimeoutPerRequest: time.Second,
    }
    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    kapi := client.NewKeysAPI(c)
    m := make(map[string]string)
    so := &client.SetOptions{}
    for i := 0; i < size; i++ {
        key := "key/x" + strconv.Itoa(i)
        keyFull := path.Join(root, key)
        value := "value-" + strconv.Itoa(i)
        kapi.Set(context.Background(), keyFull, value, so)
        //fmt.Println(res, err)
        k := strings.Replace(key, "/", ".", -1)
        //fmt.Println(key, k, value)
        m[k] = value
    }

    return m

}
