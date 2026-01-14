package etcd

import (
	"context"
	"github.com/coreos/etcd/client"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestEtcdKeyValueConfigSource(t *testing.T) {

	address := testEtcd.Address

	root := "/config101/test/kvdemo1"
	size := 10
	m := initEtcdData(address, root, size)
	c := NewEtcdKeyValueConfigSource(address, root)
	Convey("etcd kv api2", t, func() {
		keys := c.Keys()
		So(len(keys), ShouldEqual, size)
		for _, key := range keys {
			v, ok := m[key]
			//fmt.Println(key)
			v1, err := c.Get(key)
			So(ok, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, v1)
		}
	})

}

func initEtcdData(address, root string, size int) map[string]string {

	cfg := client.Config{
		Endpoints: []string{address},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)
	m := make(map[string]string)
	so := &client.SetOptions{}
	for i := 0; i < size; i++ {
		key := "key/x" + strconv.Itoa(i)
		keyFull := filepath.Join(root, key)
		value := "value-" + strconv.Itoa(i)
		kapi.Set(context.Background(), keyFull, value, so)
		//fmt.Println(res, err)
		k := strings.Replace(key, "/", ".", -1)
		//fmt.Println(key, k, value)
		m[k] = value
	}

	return m

}
