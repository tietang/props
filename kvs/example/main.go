package main

import (
	"fmt"
	"github.com/tietang/props/v3/kvs"
)

func main() {
	var conf0 kvs.ConfigSource
	var conf1 kvs.ConfigSource
	var conf2 kvs.ConfigSource
	conf := kvs.NewDefaultCompositeConfigSource(conf0, conf1, conf2)
	fmt.Println(conf.Get("x"))
}
