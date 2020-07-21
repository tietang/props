package zk

import (
	"github.com/prometheus/common/log"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/tietang/props/v3/ini"
	"github.com/tietang/props/v3/kvs"
	"github.com/tietang/props/v3/yam"
	"path"
	"strings"
)

/*
//通过key/ini_props, key所谓section，value为props格式内容，类似ini文件格式
//配置为ContentIniProps模式
// key作为section，value为props格式内容，类似ini文件格式
// key为实际key的prefix，会添加到前面
// 比如 root=configs/dev/app
// zookeeper
// 		full key=configs/dev/app/mysql
// 		value=(x1=0 \n x2=1)
// 实际key/value为： mysql.x1=0 mysql.x2=1
//ContentProps,ContentYamlContentIni 模式时，
// 其 key无实际配置意义，只作为配置分组标识，
// 值为对应的内容格式类型，读取时会将对应的内容转换这种类型
// 可以通过key后缀来标识格式类型，默认按照properties来读取

例如：
context: /configs/dev/app
如果是ContentType=ini_props，则zk nodes:

/configs/dev/app
    - /jdbc:
       ```
        url=tcp(127.0.0.1:3306)/Test?charset=utf8
        username=root
        password=root
        timeout=6s

        ```
    - /redis:
       ```
        host=192.168.1.123
        port=6379
        database=2
        timeout=6s
        password=password
        ```

如果是如下中的一个，
ContentProperties  contentType = "properties"
ContentProps       contentType = "props" //properties 别名
ContentYaml        contentType = "yaml"
ContentYam         contentType = "yam" //yaml 别名
ContentYml         contentType = "yml" //yaml 别名
ContentIni         contentType = "ini"


则zk nodes,last key以对应的配置格式，就如同文件名一样，
比如key为mysql，配置格式为yml，则在zookeeper中配置的key为mysql.yml
那么完成的zk  path就为，/configs/dev/app/mysql.yml=`k/v...`

比如properties，或者props

/configs/dev/app
    - /jdbc.props:
       ```
        jdbc.url=tcp(127.0.0.1:3306)/Test?charset=utf8
        jdbc.username=root
        jdbc.password=root
        jdbc.timeout=6s

        ```
    - /redis.props:
       ```
        redis.host=192.168.1.123
        redis.port=6379
        redis.database=2
        redis.timeout=6s
        redis.password=password
        ```

*/
type ZookeeperConfigSource struct {
	ZookeeperSource
	ContentType kvs.ContentType
}

func NewZookeeperConfigSource(watched bool, context string, conn *zk.Conn) *ZookeeperConfigSource {
	return NewZookeeperConfigSourceByName("zk:"+context, watched, context, conn, kvs.ContentAuto)
}
func NewZookeeperConfigSourceByName(name string, watched bool, context string, conn *zk.Conn, contentType kvs.ContentType) *ZookeeperConfigSource {
	s := &ZookeeperConfigSource{}
	s.ContentType = contentType
	s.name = name
	s.Watched = watched
	s.Values = make(map[string]string)
	s.conn = conn
	s.context = context
	s.init()
	return s
}

func (s *ZookeeperConfigSource) init() {
	s.findChildProperties(s.context, nil)
}

func (s *ZookeeperConfigSource) Close() {
	s.conn.Close()
}

func (s *ZookeeperConfigSource) findChildProperties(parentPath string, children []string) {
	if s.ContentType == kvs.ContentKV {
		s.findProperties(s.Watched, parentPath, children)
		return
	}
	if s.Watched && strings.HasSuffix(parentPath, DEFAULT_WATCH_KEY) {
		log.Debug("WatchNodeDataChange: ", parentPath)
		s.watchGet(parentPath)
	}

	if len(children) == 0 {
		children = s.getChildren(parentPath)
	}
	if len(children) == 0 {
		return
	}

	for _, p := range children {
		fp := path.Join(parentPath, p)
		content, err := s.getPropertiesValue(fp)
		if s.Watched && strings.HasSuffix(fp, DEFAULT_WATCH_KEY) {
			log.Debug("WatchNodeDataChange: ", fp)
			s.watchGet(fp)
		}
		if err != nil {
			continue
		}
		var ctype kvs.ContentType
		if s.ContentType == kvs.ContentAuto {
			key := path.Base(p)
			idx := strings.LastIndex(key, ".")
			if idx == -1 || idx == len(key)-1 {
				//如果获取不到格式类型，就在内容第一行注释中获取
				contentType := kvs.ReadContentType(content)
				//如果为普通文本类型，那么就默认为ContentProps
				if contentType == kvs.TextContentType {
					ctype = kvs.ContentProps
				} else {
					ctype = contentType
				}
			} else {
				ctype = kvs.ContentType(key[idx+1:])
			}
		} else {
			ctype = s.ContentType
		}

		if ctype == kvs.ContentProps || ctype == kvs.ContentProperties {
			s.findProps(content)
		} else if ctype == kvs.ContentIniProps {
			s.findIniProps(p, content)
		} else if ctype == kvs.ContentIni {
			s.findIni(content)
		} else if ctype == kvs.ContentYaml || ctype == kvs.ContentYam || ctype == kvs.ContentYml {
			s.findYaml(content)
		} else {
			log.Warn("Unsupported format：", s.ContentType)
		}

	}

}
func (s *ZookeeperConfigSource) findYaml(content string) {
	props := yam.ByYaml(content)
	if props != nil {
		s.SetAll(props.Values)
	}
}

func (s *ZookeeperConfigSource) findIni(content string) {
	props := ini.ByIni(content)
	if props != nil {
		s.SetAll(props.Values)
	}

}

func (s *ZookeeperConfigSource) findProps(content string) {
	props := kvs.ByProperties(content)
	if props != nil {
		s.SetAll(props.Values)
	}
}
func (s *ZookeeperConfigSource) findIniProps(key, content string) {
	props := kvs.ByProperties(content)
	if props != nil {
		prefix := path.Base(key)
		for key, value := range props.Values {
			k := prefix + "." + key
			s.Set(k, value)
		}
	}

}
