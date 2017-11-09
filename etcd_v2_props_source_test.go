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

func TestEtcdPropsConfigSource(t *testing.T) {

    address := testEtcd.Address

    root := "/config101/props/demo/v2"
    size := 10
    inilen := 3
    m := initEtcdV2PropsData(address, root, size, inilen)
    c := NewEtcdPropsConfigSource(address, root)
    Convey("etcd props api2", t, func() {
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

func initEtcdV2PropsData(address, root string, size, inilen int) map[string]string {

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
        key := "key-" + strconv.Itoa(i)
        keyFull := path.Join(root, key)
        value := ""

        for j := 0; j < inilen; j++ {
            kk := key + "." + "x" + strconv.Itoa(j)
            val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
            value += "x" + strconv.Itoa(j) + "=" + val + "\n"
            k := strings.Replace(kk, "/", ".", -1)
            //fmt.Println(key, k, value)
            m[k] = val
        }
        kapi.Set(context.Background(), keyFull, value, so)
    }

    return m

}
