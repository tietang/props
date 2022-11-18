package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/consul"
	"github.com/tietang/props/v3/kvs"
	"time"
)

/*
*

conf/ini/demo1 = `

[x0]
y0=val0
y1=val1

[x1]
y0=val0
y1=val1

`

conf/ini_props/ =
conf/ini_props/x0 = y0=val0
y1=val1
conf/ini_props/x1 = y1=val1

conf/kv/ =
conf/kv/x0/y0 = val0
conf/kv/x1/y1 = val1

conf/props/ =
conf/props/demo1 = x0.y0=val0
x1.y1=val1
conf/props/demo2 = k.y=val0
k2.y2=val2
*/
func main_3() {
	address := "127.0.0.1:8500"

	root := "conf/"
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	kv := client.KV()
	q := &api.QueryOptions{}

	keys, _, err := kv.Keys(root, "", q)
	if err != nil {
		log.Error(err)
		return
	}
	for _, k := range keys {
		kv, _, err := kv.Get(k, q)
		if err != nil {
			log.Error(err)
			continue
		}
		value := string(kv.Value)
		fmt.Println(k, "=", value)
	}

	fmt.Println("\n kv:")
	root = "conf/kv/"
	conf := consul.NewConsulConfigSourceByName("zk-"+root, address, root, kvs.ContentKV, time.Second*10)
	keys = conf.Keys()
	for _, k := range keys {
		v := conf.GetDefault(k, "")
		fmt.Println(k, "=", v)
	}
	fmt.Println("--------------")
	fmt.Println("\n ini_props:")
	root = "conf/ini_props/"
	conf = consul.NewConsulConfigSourceByName("zk-"+root, address, root, kvs.ContentIniProps, time.Second*10)
	keys = conf.Keys()
	for _, k := range keys {
		v := conf.GetDefault(k, "")
		fmt.Println(k, "=", v)
	}
	fmt.Println("--------------")

	fmt.Println("\n props:")
	root = "conf/props/"
	conf = consul.NewConsulConfigSourceByName("zk-"+root, address, root, kvs.ContentProps, time.Second*10)
	keys = conf.Keys()
	for _, k := range keys {
		v := conf.GetDefault(k, "")
		fmt.Println(k, "=", v)
	}
	fmt.Println("--------------")

	fmt.Println("\n ini:")
	root = "conf/ini/"
	conf = consul.NewConsulConfigSourceByName("zk-"+root, address, root, kvs.ContentIni, time.Second*10)
	keys = conf.Keys()
	for _, k := range keys {
		v := conf.GetDefault(k, "")
		fmt.Println(k, "=", v)
	}
	fmt.Println("--------------")

}
