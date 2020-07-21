package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	s := time.Second * 10
	v := reflect.ValueOf(s)
	fmt.Println(v.Kind())
	fmt.Println(v.Kind())
}
