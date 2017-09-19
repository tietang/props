package main

import (
    "github.com/tietang/props"
    "fmt"
    "time"
)

func main3() {

    root := "config/app1/dev"
    address := "127.0.0.1:8500"

    p := props.NewConsulPropsConfigSourceByName("consul-props", address, root, 10*time.Second)
    fmt.Println(p.Get(""))
}
