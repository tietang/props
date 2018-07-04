package zk

import (
    "testing"
    "time"
    . "github.com/smartystreets/goconvey/convey"
    //"github.com/lunny/log"
    "strconv"
    "github.com/samuel/go-zookeeper/zk"
    "fmt"
    "strings"
    "path"
)

func initIniData() (*ZookeeperPropsConfigSource, map[string]string) {
    size := 10
    inilen := 3
    //urls:=[]string{"172.16.1.248:2181"}
    urls := []string{"127.0.0.1:2181"}
    c, ch, err := zk.Connect(urls, 20*time.Second)
    conn := c
    if err != nil {
        panic(err)
    }
    event := <-ch

    fmt.Println(event)
    root := "/config_ini/app1/dev"
    //fmt.Println("d:  ", conn.State().String(), err, contexts[0])
    kv := initZkIniData(conn, root, size, inilen)
    zics := NewZookeeperPropsConfigSource("zookeeper-props", root, conn)
    return zics, kv
}

func initZkIniData(conn *zk.Conn, root string, size, inilen int) map[string]string {
    kv := make(map[string]string)
    for i := 0; i < size; i++ {
        key := "key-" + strconv.Itoa(i)
        keyPath := path.Join(root, key)

        value := ""

        for j := 0; j < inilen; j++ {
            kk := key + "." + "x" + strconv.Itoa(j)
            val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
            value += "x" + strconv.Itoa(j) + "=" + val + "\n"
            k := strings.Replace(kk, "/", ".", -1)
            //fmt.Println(key, k, value)
            kv[k] = val
        }

        if !ZkExits(conn, keyPath) {
            _, err := ZkCreateString(conn, keyPath, value)
            if err == nil {
                //log.Println(path)
            }
            //log.Println(err)
        }
    }
    return kv

    //fmt.Println(len(kv))
}

func TestReadZkIni(t *testing.T) {
    zics, kv := initIniData()
    Convey("Get", t, func() {

        keys := zics.Keys()
        //fmt.Println(len(kv), len(keys))
        //fmt.Println(kv)
        //fmt.Println(keys)
        So(len(keys), ShouldBeGreaterThanOrEqualTo, len(kv))
        Convey("验证", func() {
            for _, k := range keys {
                v1, _ := zics.Get(k)
                v2 := kv[k]
                So(v1, ShouldEqual, v2)

                //fmt.Println(k, "=", v1, "   ")
            }
        })

    })

    //conn.Close()

}
