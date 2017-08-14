package main

import (
    "github.com/tietang/props"
    "github.com/samuel/go-zookeeper/zk"
    "time"
    "fmt"
)

func main() {
    context := "/configurations/demo/dev/app1"

    c, ch, err := zk.Connect([]string{"172.16.1.248:2181"}, 2*time.Second)
    if err != nil {
        panic(err)
    }
    for {
        event := <-ch
        fmt.Println(event)
        if event.State == zk.StateConnected {
            break
        }

    }

    z := props.NewZookeeperConfigSource("zk", context, c)
    v, err := z.Get("app.xx.key-0")
    fmt.Println(v)
    fmt.Println(err)
    z.Watch("app/xx", func(children []string, event zk.Event) {
        fmt.Println("Watch:   ", event, len(z.Keys()))
        for _, key := range z.Keys() {
            fmt.Println("new value:  ", z.GetDefault(key, key))

        }
    })
    z.WatchChildren("app/xx", func(children []string, event zk.Event) {

        fmt.Println("WatchChildren:   ", event, len(z.Keys()))
        for _, key := range z.Keys() {
            fmt.Println("new value:  ", z.GetDefault(key, key))
        }
    })
    for {
        event := <-ch
        fmt.Println(event)

    }
}
