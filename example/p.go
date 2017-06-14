package main

import (
	"github.com/tietang/props"
	"fmt"
)

func main() {
	p, err := props.ReadPropertyFile("config.properties")
	if err != nil {
		panic(err)
	}
	stringValue, err := p.Get("prefix.key1")
	fmt.Println(stringValue, err)
	//如果不存在，则返回默认值
	stringDefaultValue := p.GetDefault("prefix.key1", "default value")
	fmt.Println(stringDefaultValue)
	boolValue, err := p.GetBool("prefix.key2")
	fmt.Println(boolValue)
	boolDefaultValue := p.GetBoolDefault("prefix.key2", false)
	fmt.Println(boolDefaultValue)
	intValue, err := p.GetInt("prefix.key3")
	fmt.Println(intValue)
	intDefaultValue := p.GetIntDefault("prefix.key3", 1)
	fmt.Println(intDefaultValue)
	floatValue, err := p.GetFloat64("prefix.key4")
	fmt.Println(floatValue)
	floatDefaultValue := p.GetFloat64Default("prefix.key4", 1.2)
	fmt.Println(floatDefaultValue)
}
