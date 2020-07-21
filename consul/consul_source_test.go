package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/props/v3/kvs"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestConsulIniPropsConfigSource(t *testing.T) {
	address := "127.0.0.1:8500"
	//address := "172.16.1.248:8500"
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	kv := client.KV()

	root := "config101/test/inidemo"
	size := 10
	inilen := 3

	Convey("consul kv", t, func() {
		m := initPropsConsulData(address, root, size, inilen)
		c := NewConsulConfigSource(address, root)
		keys := c.Keys()
		So(len(keys), ShouldEqual, size*inilen)
		for _, key := range keys {
			v, ok := m[key]
			//fmt.Println(key)
			v1, err := c.Get(key)
			So(ok, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, v1)
		}
	})

	Convey("Consul For Props", t, func() {
		root := "config_props/dev/app"
		kv := initConsulDataForProps(kv, root, 2, 2)
		conf := NewConsulConfigSource(address, root)
		keys := conf.Keys()
		So(len(keys), ShouldEqual, len(kv))
		Convey("验证", func() {
			for _, k := range keys {
				v1, _ := conf.Get(k)
				v2 := kv[k]
				So(v1, ShouldEqual, v2)
			}
		})

	})
	Convey("Consul for IniProps", t, func() {
		root := "config_ini_props/dev/app"
		kv := initConsulDataForIniProps(kv, root, 10, 10)
		conf := NewConsulConfigSourceByName("zk-"+root, address, root, kvs.ContentIniProps, time.Second*10)

		keys := conf.Keys()
		So(len(keys), ShouldEqual, 10*10)
		for _, key := range keys {
			v, ok := kv[key]
			//fmt.Println(key)
			v1, err := conf.Get(key)
			So(ok, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, v1)
		}
	})
	Convey("Consul For ini", t, func() {
		root := "config_ini/dev/app"
		kv := initConsulDataForIni(kv, root, 2, 2)
		conf := NewConsulConfigSource(address, root)
		keys := conf.Keys()

		So(len(keys), ShouldEqual, len(kv))

		Convey("验证", func() {
			for _, k := range keys {
				v1, _ := conf.Get(k)
				v2 := kv[k]
				So(v1, ShouldEqual, v2)
			}
		})

	})

}

func TestConsulConfigSourceForKV(t *testing.T) {

	address := testConsul.Address

	root := "config_kv/test/kv_demo1"
	size := 2
	m := initConsulDataForKV(address, root, size)
	conf := NewConsulConfigSourceByName("zk-"+root, address, root, kvs.ContentKV, time.Second*10)

	Convey("consul kv", t, func() {
		keys := conf.Keys()
		So(len(keys), ShouldEqual, size)
		for _, key := range keys {
			v, ok := m[key]
			v1, err := conf.Get(key)
			So(ok, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, v1)
		}
	})

}

func initConsulDataForKV(address, root string, size int) map[string]string {
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	m := make(map[string]string)
	kv := client.KV()
	wq := &api.WriteOptions{}
	for i := 0; i < size; i++ {
		key := "key/x" + strconv.Itoa(i)
		keyFull := path.Join(root, key)
		value := "value-" + strconv.Itoa(i)
		kvp := &api.KVPair{
			Key:   keyFull,
			Value: []byte(value),
		}
		kv.Put(kvp, wq)
		k := strings.Replace(key, "/", ".", -1)
		//fmt.Println(key, k, value)
		m[k] = value
	}

	return m

}

func initPropsConsulData(address, root string, size, len int) map[string]string {
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	m := make(map[string]string)
	kv := client.KV()
	wq := &api.WriteOptions{}

	for i := 0; i < size; i++ {
		key := "key-" + strconv.Itoa(i)
		keyFull := path.Join(root, key)

		value := ""

		for j := 0; j < len; j++ {
			val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
			pkey := "x" + strconv.Itoa(i) + "-y" + strconv.Itoa(j)
			value += pkey + "=" + val + "\n"
			m[pkey] = val
		}
		kvp := &api.KVPair{
			Key:   keyFull,
			Value: []byte(value),
		}
		_, err = kv.Put(kvp, wq)
		if err != nil {
			logrus.Error(err)
		}
	}

	return m

}

func initConsulDataForProps(kv *api.KV, root string, size, inilen int) map[string]string {
	w := &api.WriteOptions{}
	_, err := kv.Delete(root, w)
	if err != nil {
		logrus.Error(err)
	}
	m := make(map[string]string)
	for i := 0; i < size; i++ {
		key := "key" + strconv.Itoa(i) + ".props"
		keyPath := path.Join(root, key)

		value := ""

		for j := 0; j < inilen; j++ {
			val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
			pkey := "x" + strconv.Itoa(i) + "-y" + strconv.Itoa(j)
			value += pkey + "=" + val + "\n"
			//fmt.Println(key, k, value)
			m[pkey] = val
		}

		_, err = kv.Delete(keyPath, w)
		if err != nil {
			logrus.Error(err)
		}

		kvp := &api.KVPair{
			Key:   keyPath,
			Value: []byte(value),
		}
		_, err = kv.Put(kvp, w)
		if err != nil {
			logrus.Error(err)
		}
	}
	return m

}
func initConsulDataForIni(kv *api.KV, root string, size, inilen int) map[string]string {
	w := &api.WriteOptions{}
	_, err := kv.Delete(root, w)
	if err != nil {
		logrus.Error(err)
	}
	m := make(map[string]string)
	keyPath := path.Join(root, "test.ini")

	value := ""
	for i := 0; i < size; i++ {
		section := "key" + strconv.Itoa(i)

		value += "\n[" + section + "]\n"

		for j := 0; j < inilen; j++ {
			val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
			pkey := "x" + strconv.Itoa(i) + "-y" + strconv.Itoa(j)
			value += pkey + "=" + val + "\n"
			kk := section + "." + pkey
			m[kk] = val
		}

	}
	_, err = kv.Delete(keyPath, w)
	if err != nil {
		logrus.Error(err)
	}

	kvp := &api.KVPair{
		Key:   keyPath,
		Value: []byte(value),
	}
	_, err = kv.Put(kvp, w)
	if err != nil {
		logrus.Error(err)
	}

	return m
}

func initConsulDataForIniProps(kv *api.KV, root string, size, inilen int) map[string]string {
	w := &api.WriteOptions{}
	_, err := kv.Delete(root, w)
	if err != nil {
		logrus.Error(err)
	}
	m := make(map[string]string)
	for i := 0; i < size; i++ {
		key := "key-" + strconv.Itoa(i)
		keyPath := path.Join(root, key)

		value := ""

		for j := 0; j < inilen; j++ {
			kk := key + "." + "x" + strconv.Itoa(j)
			val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
			value += "x" + strconv.Itoa(j) + "=" + val + "\n"
			k := strings.Replace(kk, "/", ".", -1)
			m[k] = val
		}
		_, err = kv.Delete(keyPath, w)
		if err != nil {
			logrus.Error(err)
		}

		kvp := &api.KVPair{
			Key:   keyPath,
			Value: []byte(value),
		}
		_, err = kv.Put(kvp, w)
		if err != nil {
			logrus.Error(err)
		}
	}

	return m

}
