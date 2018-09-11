package main

import (
    "github.com/tietang/props/zk"
    "time"
    "fmt"
)

func main() {
    contexts := []string{
        "/config/root",
        "/config/app1",
    }

    conf := zk.NewZookeeperCompositeConfigSource(contexts, []string{"172.16.1.248:2181"}, 2*time.Second)

    fmt.Println(conf.Keys())
    fmt.Println(conf.GetDefault("key1","__def_value"))
    fmt.Println(conf.GetDefault("foo","__def_value"))
    fmt.Println(conf.GetDefault("password","__def_value"))
    fmt.Println(conf.GetDefault("redis_url","__def_value"))
    fmt.Println(conf.GetDefault("username","__def_value"))
}
