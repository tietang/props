package main

import (
	"fmt"

	"github.com/tietang/props/v3/kvs"
	"github.com/tietang/props/v3/tom"
)

func main() {
	file := kvs.GetCurrentFilePath("Cargo.toml", 2)
	props, err := tom.ReadTomlFile(file)
	if err != nil {
		panic(err)
	}
	for k, v := range props.Values {
		fmt.Println(k, "=", v)
	}
}
