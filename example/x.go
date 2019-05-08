package main

import (
	"fmt"
	"github.com/tietang/props/kvs"
	"path"
	"strings"
)

func main_1() {
	var ctype kvs.ContentType

	k := "configs/dev/resk/mysql."
	key := path.Base(k)
	idx := strings.LastIndex(key, ".")
	fmt.Println(idx)
	if idx == -1 || idx == len(key)-1 {
		ctype = kvs.ContentProps
	} else {
		ctype = kvs.ContentType(key[idx+1:])
	}

	fmt.Println(ctype)

}
