package main

import (
	"fmt"

	"github.com/tietang/props/v3/kvs"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini")
	fmt.Println(file)
	fmt.Println(kvs.CurrentFile())
}
