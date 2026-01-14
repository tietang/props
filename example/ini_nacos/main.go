package main

import (
	"fmt"
	"github.com/tietang/props/v3/nacos"
)

func main() {
	address := "10.99.71.54:8848"
	conf := nacos.NewNacosClientCompositeConfigSource(address, "dev", "dzpl", []string{"collector", "collector-go2"})
	fmt.Println(conf.Get("app.server.port"))

}
