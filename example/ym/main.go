package main

import (
	"fmt"

	"github.com/tietang/props/v3/kvs"
	"github.com/tietang/props/v3/yam"
)

func main() {
	file := kvs.GetCurrentFilePath("boot.yml")
	props, err := yam.ReadYamlFile(file)
	if err != nil {
		panic(err)
	}
	for k, v := range props.Values {
		fmt.Println(k, "=", v)
	}
}
