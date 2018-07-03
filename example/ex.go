package main

import (
	"fmt"
	"github.com/tietang/props/kvs"
)

func main2() {
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
	var pcs kvs.ConfigSource
	var err error
	//通过文件名，文件名作为ConfigSource name
	pcs = kvs.NewPropertiesConfigSource("config.properties")

	stringValue, err := pcs.Get("prefix.key1")
	fmt.Println(stringValue, err)
	//如果不存在，则返回默认值
	stringDefaultValue := pcs.GetDefault("prefix.key1", "default value")
	fmt.Println(stringDefaultValue)
	boolValue, err := pcs.GetBool("prefix.key2")
	fmt.Println(boolValue)
	boolDefaultValue := pcs.GetBoolDefault("prefix.key2", false)
	fmt.Println(boolDefaultValue)
	intValue, err := pcs.GetInt("prefix.key3")
	fmt.Println(intValue)
	intDefaultValue := pcs.GetIntDefault("prefix.key3", 1)
	fmt.Println(intDefaultValue)
	floatValue, err := pcs.GetFloat64("prefix.key4")
	fmt.Println(floatValue)
	floatDefaultValue := pcs.GetFloat64Default("prefix.key4", 1.2)
	fmt.Println(floatDefaultValue)
	//

	//

}
