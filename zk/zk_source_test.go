package zk

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/props/v3/kvs"
	"path"
	"strings"
	"testing"
	//"github.com/lunny/log"
	"strconv"
	"time"
)

type logWriter struct {
	t *testing.T
	p string
}

func (lw logWriter) Write(b []byte) (int, error) {
	lw.t.Logf("%s%s", lw.p, string(b))
	return len(b), nil
}

func TestReadZkForProps(t *testing.T) {

	//urls:=[]string{"172.16.1.248:2181"}

	urls := []string{"127.0.0.1:2181"}
	c, ch, err := zk.Connect(urls, 20*time.Second)
	if err != nil {
		panic(err)
	}
	event := <-ch

	fmt.Println(event)

	Convey("Zookeeper For Props", t, func() {
		root := "/config_props/dev/app"
		kv := initZkDataForProps(c, root, 2, 2)
		conf := NewZookeeperConfigSource(true, root, c)
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

	Convey("Zookeeper for IniProps", t, func() {
		root := "/config_ini_props/dev/app"
		kv := initZkDataForIniProps(c, root, 10, 10)
		conf := NewZookeeperConfigSourceByName("zk-"+root, true, root, c, kvs.ContentIniProps)

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
	Convey("Zookeeper For ini", t, func() {
		root := "/config_ini/dev/app"
		kv := initZkDataForIni(c, root, 2, 2)
		conf := NewZookeeperConfigSource(true, root, c)
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
	//conn.Close()

}

func TestReadZkForKv(t *testing.T) {

	//urls:=[]string{"172.16.1.248:2181"}
	urls := []string{"127.0.0.1:2181"}
	c, ch, err := zk.Connect(urls, 20*time.Second)
	if err != nil {
		panic(err)
	}
	event := <-ch

	fmt.Println(event)
	root := "/config_kv/app1/dev"
	kv := initZkDataForKV(c, root, 10)
	conf := NewZookeeperConfigSourceByName("zk-"+root, true, root, c, kvs.ContentKV)

	Convey("Get", t, func() {

		keys := conf.Keys()
		//fmt.Println(len(kv), len(keys))
		//fmt.Println(kv)
		//fmt.Println(keys)
		So(len(keys), ShouldBeGreaterThanOrEqualTo, len(kv))
		Convey("验证", func() {
			for _, k := range keys {
				v1, _ := conf.Get(k)
				v2 := kv[k]
				So(v1, ShouldEqual, v2)

				//fmt.Println(k, "=", v1, "   ")
			}
		})

	})

	//conn.Close()

}

func initZkDataForKV(conn *zk.Conn, root string, size int) map[string]string {
	kv := make(map[string]string)
	for i := 0; i < size; i++ {
		key := "key-" + strconv.Itoa(i)
		value := "value-" + strconv.Itoa(i)
		p := root + "/app/xx/" + key
		vkey := "app.xx." + key

		s, ok := ZkExists(conn, p)
		if ok {
			err := conn.Delete(p, s.Version)
			fmt.Println(err)
		} else {
			kv[vkey] = value
		}
		_, err := ZkCreateString(conn, p, value)
		if err == nil {
			//log.Println(path)
			kv[vkey] = value
		}

	}
	return kv

	//fmt.Println(len(kv))
}

func initZkDataForProps(conn *zk.Conn, root string, size, inilen int) map[string]string {
	err := ZkDelete(conn, root)
	if err != nil {
		logrus.Error(err)
	}
	kv := make(map[string]string)
	for i := 0; i < size; i++ {
		key := "key" + strconv.Itoa(i) + ".props"
		keyPath := path.Join(root, key)

		value := ""

		for j := 0; j < inilen; j++ {
			val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
			pkey := "x" + strconv.Itoa(i) + "-y" + strconv.Itoa(j)
			value += pkey + "=" + val + "\n"
			//fmt.Println(key, k, value)
			kv[pkey] = val
		}

		err := ZkDelete(conn, keyPath)
		if err != nil {
			logrus.Error(err)
		}

		_, err = ZkCreateString(conn, keyPath, value)
		if err == nil {
			logrus.Error(err)
		}
	}
	return kv

	//fmt.Println(len(kv))
}

func initZkDataForIni(conn *zk.Conn, root string, size, inilen int) map[string]string {
	err := ZkDelete(conn, root)
	if err != nil {
		logrus.Error(err)
	}
	kv := make(map[string]string)
	keyPath := path.Join(root, "test.ini")

	value := ""
	for i := 0; i < size; i++ {
		section := "key" + strconv.Itoa(i)

		value += "\n[" + section + "]\n"

		for j := 0; j < inilen; j++ {
			val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
			pkey := "x" + strconv.Itoa(i) + "-y" + strconv.Itoa(j)
			value += pkey + "=" + val + "\n"
			//fmt.Println(key, k, value)
			kk := section + "." + pkey
			kv[kk] = val
		}

	}
	err = ZkDelete(conn, keyPath)
	if err != nil {
		logrus.Error(err)
	}
	_, err = ZkCreateString(conn, keyPath, value)
	if err == nil {
		//log.Println(path)
	}

	return kv

	//fmt.Println(len(kv))
}

func initZkDataForIniProps(conn *zk.Conn, root string, size, inilen int) map[string]string {
	err := ZkDelete(conn, root)
	if err != nil {
		logrus.Error(err)
	}
	kv := make(map[string]string)
	for i := 0; i < size; i++ {
		key := "key-" + strconv.Itoa(i)
		keyPath := path.Join(root, key)

		value := ""

		for j := 0; j < inilen; j++ {
			kk := key + "." + "x" + strconv.Itoa(j)
			val := "value-" + strconv.Itoa(i) + strconv.Itoa(j)
			value += "x" + strconv.Itoa(j) + "=" + val + "\n"
			k := strings.Replace(kk, "/", ".", -1)
			//fmt.Println(key, k, value)
			kv[k] = val
		}

		err := ZkDelete(conn, keyPath)
		if err != nil {
			logrus.Error(err)
		}
		_, err = ZkCreateString(conn, keyPath, value)
		if err == nil {
			//log.Println(path)
		}
	}
	return kv

	//fmt.Println(len(kv))
}
