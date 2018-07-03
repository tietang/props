package main

import (
	"time"
	"fmt"
	"github.com/tietang/props/kvs"
	"github.com/tietang/props/zk"
)

func main() {
	//p, err := props.ReadPropertyFile("config.properties")
	//if err != nil {
	//	panic(err)
	//}
	//stringValue, err := p.Get("prefix.key1")
	////如果不存在，则返回默认值
	//stringDefaultValue := p.GetDefault("prefix.key1", "default value")
	//boolValue, err := p.GetBool("prefix.key1")
	//boolDefaultValue := p.GetBoolDefault("prefix.key1", false)
	//intValue, err := p.GetInt("prefix.key1")
	//intDefaultValue := p.GetIntDefault("prefix.key1", 1)
	//floatValue, err := p.GetFloat64("prefix.key1")
	//floatDefaultValue := p.GetFloat64Default("prefix.key1", 1.2)
	//
	////
	//通过文件名，文件名作为ConfigSource name
	pcs1 := kvs.NewPropertiesConfigSource("config.properties")
	//指定名称和文件名
	pcs2 := kvs.NewPropertiesConfigSourceByFile("config", "config.properties")

	urls := []string{"172.16.1.248:2181"}
	contexts := []string{"/configs/apps", "/configs/users"}
	zccs := zk.NewZookeeperCompositeConfigSource(contexts, urls, time.Second*3)
	configSources := []kvs.ConfigSource{pcs1, pcs2, zccs, }
	ccs := kvs.NewDefaultCompositeConfigSource(configSources...)

	//

	stringValue, err := ccs.Get("prefix.key1")
	fmt.Println(stringValue, err)
	//如果不存在，则返回默认值
	stringDefaultValue := ccs.GetDefault("prefix.key1", "default value")
	fmt.Println(stringDefaultValue)
	boolValue, err := ccs.GetBool("prefix.key2")
	fmt.Println(boolValue)
	boolDefaultValue := ccs.GetBoolDefault("prefix.key2", false)
	fmt.Println(boolDefaultValue)
	intValue, err := ccs.GetInt("prefix.key3")
	fmt.Println(intValue)
	intDefaultValue := ccs.GetIntDefault("prefix.key3", 1)
	fmt.Println(intDefaultValue)
	floatValue, err := ccs.GetFloat64("prefix.key4")
	fmt.Println(floatValue)
	floatDefaultValue := ccs.GetFloat64Default("prefix.key4", 1.2)
	fmt.Println(floatDefaultValue)
}
