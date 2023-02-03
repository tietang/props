package main

import (
	"fmt"
	"github.com/tietang/props/v3/apollo"
)

func main() {
	////http://106.12.25.204:8080/configfiles/json/we/default/application?ip=1.1.1.1
	//
	//http://81.68.181.139/configfiles/json/SampleApp/default/application?ip=1.1.1.1
	//{config_server_url}/configfiles/json/{appId}/{clusterName}/{namespaceName}?ip={clientIp}
	a := apollo.NewApolloConfigSourceWithSecret("81.68.181.139:8080", "SampleApp", "21ec500b3bf34bbd8e4319ff3c22e301", []string{
		"application", "brian.yaml", "test_ini.txt", "acewan.properties",
	})
	a.AddChangeListener("timeout", func(k, v string) {
		fmt.Println(k, v)
	})
	a.AddChangeListener("key_inter_2", func(k, v string) {
		fmt.Println(k, v)
	})
	a.AddWatchNamespaces("application")

	keys := a.Keys()
	for _, key := range keys {
		value := a.GetDefault(key, "null")
		fmt.Println(key, "=", value)
	}

	select {}

}
