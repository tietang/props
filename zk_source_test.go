package props

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    //"github.com/lunny/log"
    "strconv"
    "github.com/samuel/go-zookeeper/zk"
    "time"
    "fmt"
)

type logWriter struct {
    t *testing.T
    p string
}

func (lw logWriter) Write(b []byte) (int, error) {
    lw.t.Logf("%s%s", lw.p, string(b))
    return len(b), nil
}


func TestReadZk(t *testing.T) {

    //urls:=[]string{"172.16.1.248:2181"}
    urls := []string{"127.0.0.1:2181"}
    c, ch, err := zk.Connect(urls, 20*time.Second)
    zkConn := c
    if err != nil {
        panic(err)
    }
    event := <-ch

    fmt.Println(event)
    zs, kv := initData(zkConn)

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

func initData(conn *zk.Conn) (*ZookeeperConfigSource, map[string]string) {

    size := 10

    root := "/config_kv/app1/dev"
    //fmt.Println("d:  ", conn.State().String(), err, contexts[0])
    kv := initZkData(conn, root, size)
    zics := NewZookeeperConfigSource("zookeeper-props", root, conn)
    return zics, kv
}

func initZkData(conn *zk.Conn, root string, size int) map[string]string {
    kv := make(map[string]string)
    for i := 0; i < size; i++ {
        key := "key-" + strconv.Itoa(i)
        value := "value-" + strconv.Itoa(i)
        p := root + "/app/xx/" + key
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
    return kv

    //fmt.Println(len(kv))
}
