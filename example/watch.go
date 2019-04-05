package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"net/http"
	"time"
)

func main() {

	root := "configs/dev/resk/"

	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	config.WaitTime = time.Minute * 1
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	kv := client.KV()
	q := &api.QueryOptions{
		WaitIndex: uint64(time.Now().Unix()),
		WaitTime:  time.Second * 10,
	}
	kvals, qm, err := kv.List(root, q)
	fmt.Println(qm)
	fmt.Println(err)
	fmt.Printf("%+v", kvals)
	http.HandleFunc("/notice", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	http.ListenAndServe(":18081", nil)

}
