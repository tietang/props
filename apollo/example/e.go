package main

import (
	"fmt"
	"github.com/tietang/props/v3/apollo"
)

func main() {
	//http://106.12.25.204:8080/configfiles/json/we/default/application?ip=1.1.1.1
	a := apollo.NewApolloConfigSource("106.12.25.204:8080", "we", []string{
		"application", "development.redis",
	})
	keys := a.Keys()
	for _, key := range keys {
		value := a.GetDefault(key, "null")
		fmt.Println(key, "=", value)
	}

}
