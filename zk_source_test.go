package props

import (
    "testing"
    "time"
    . "github.com/smartystreets/goconvey/convey"
    //"github.com/lunny/log"
    "strconv"
    "github.com/samuel/go-zookeeper/zk"
    "fmt"
)

var zs *CompositeConfigSource
var contexts []string
var keyPrefix string
var kv map[string]string
var conn *zk.Conn
var size int

func init() {
    size = 10
    contexts = []string{"/configurations/demo/dev/app1", "/configurations/demo/dev/apps"}
    keyPrefix = "app.xx."

    kv = make(map[string]string)
    c, ch, err := zk.Connect([]string{"172.16.1.248:2181"}, 2 * time.Second)
    conn = c
    if err != nil {
        panic(err)
    }
    event := <-ch

    fmt.Println(event)
    //fmt.Println("d:  ", conn.State().String(), err, contexts[0])
    initZkData()
    zs = NewZookeeperCompositeConfigSourceByConn(contexts, conn)
}

func initZkData() {
    for i := 0; i < size; i++ {
        key := "key-" + strconv.Itoa(i)
        value := "value-" + strconv.Itoa(i)
        p := contexts[0] + "/app/xx/" + key
        vkey := "app.xx." + key
        if !ZkExits(conn, p) {
            _, err := ZkCreateString(conn, p, value)
            if err == nil {
                //log.Println(path)
                kv[vkey] = value
            }
            //log.Println(err)
        } else {
            kv[vkey] = value
        }

    }

    //fmt.Println(len(kv))
}
func TestReadZk(t *testing.T) {

    Convey("Get", t, func() {

        keys := zs.Keys()
        //fmt.Println(len(kv), len(keys))
        //fmt.Println(kv)
        //fmt.Println(keys)
        So(len(keys), ShouldBeGreaterThanOrEqualTo, len(kv))
        Convey("验证", func() {
            for _, k := range keys {
                v1, _ := zs.Get(k)
                v2 := kv[k]
                So(v1, ShouldEqual, v2)

                //fmt.Println(k, "=", v1, "   ")
            }
        })

    })

    //conn.Close()

}