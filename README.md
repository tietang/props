# props
 
 
统一的配置工具库，将各种配置源抽象或转换为类似properties格式的key/value，并提供统一的API来访问这些key/value。

## 支持的配置源：

- properties文件
- ini文件
- zookeeper k/v
- zookeeper k/props 
- consul k/v
- consul k/props

## key/value支持的数据类型：

- key只支持string
- value 5种数据类型的支持：
	- string
	- int
	- float64
	- bool
	- time.Duration：支持ms和s级配置,默认为ms


## Install

运行deps.sh安装依赖。

**或者**

参考 [golang dep](<https://github.com/golang/dep>)使用dep命令来安装。


## 配置源和配置形式使用方法：

### properties格式文件

格式：`[key][=|:][value] \n`
每行为key/value键值对 ,用`=`或`：`分割，key可以是除了`=`和`:`、以及空白字符的任何字符

例子：

`server.port=8080`

或者

`server.port: 8080`


### 通过props.ReadPropertyFile读取文件

```golang

p, err := props.ReadPropertyFile("config.properties")
if err != nil {
	panic(err)
}
stringValue, err := p.Get("prefix.key1")
//如果不存在，则返回默认值
stringDefaultValue := p.GetDefault("prefix.key1", "default value")
boolValue, err := p.GetBool("prefix.key1")
boolDefaultValue := p.GetBoolDefault("prefix.key1", false)
intValue, err := p.GetInt("prefix.key1")
intDefaultValue := p.GetIntDefault("prefix.key1", 1)
floatValue, err := p.GetFloat64("prefix.key1")
floatDefaultValue := p.GetFloat64Default("prefix.key1", 1.2)
v, err := p.GetDuration("k9")
v := p.GetDurationDefault("k12", 1*time.Second)

```

#### 通过props.NewProperties()从io.Reader中读取

```
 p := props.NewProperties()
 p.Load(strings.NewReader("some data"))
 p.Load(bytes.NewReader([]byte("some data")))
```

#### 通过props.NewPropertiesConfigSource()

```
file := "/path/to/config.properties"
p := props.NewPropertiesConfigSource(file)
p = props.NewPropertiesConfigSourceByFile("name", file)
//通过map构造内存型
m := make(map[string]string)
m["key"]="value"
p = props.NewPropertiesConfigSourceByMap("name", m)

```
#### Properties ConfigSource

```golang

var cs props.ConfigSource
//
cs = props.NewPropertiesConfigSource("config.properties")
cs = props.NewPropertiesConfigSourceByFile("config", "config.properties")


stringValue, err := cs.Get("prefix.key1")
//如果不存在，则返回默认值
stringDefaultValue := cs.GetDefault("prefix.key1", "default value")
boolValue, err := cs.GetBool("prefix.key2")
boolDefaultValue := cs.GetBoolDefault("prefix.key2", false)
intValue, err := cs.GetInt("prefix.key3")
intDefaultValue := cs.GetIntDefault("prefix.key3", 1)
floatValue, err := cs.GetFloat64("prefix.key4")
floatDefaultValue := cs.GetFloat64Default("prefix.key4", 1.2)

```




### ini格式文件。

格式：参考 [wiki百科：INI_file](<https://en.wikipedia.org/wiki/INI_file>)
	
```ini
[section]
[key1][=|:][value1] 
[key1][=|:][value1]
...
```

不支持sub section

例子：

```ini
[server]
port: 8080
read.timeout=6000ms

[client]
connection.timeout=6s
query.timeout=6s
```
	
#### 使用方法：
	
```golang
file := "/path/to/config.ini"
p := props.NewIniFileConfigSource(file)
p = props.NewIniFileConfigSourceByFile("name", file)
```
	
### zookeeper 

支持key/value和key/properties配置形式，key/properties配置和ini类似，将key作为section name。
key/value形式，将path去除root path部分并替换`/`为`.`作为key。
key/properties形式，在root path下读取所有子节点，将子节点名称作为section name，value为子properties格式内存，通过子节点名称和子properties中的key组合成新的key作为key。
 
#### by zookeeper key/value

##### 基本例子

```golang
root := "/config/kv/app1/dev"
var conn *zk.Conn
p := props.NewZookeeperConfigSource("zookeeper-kv", root, conn)
```

##### CompositeConfigSource多context例子

```golang
var cs props.ConfigSource
urls := []string{"172.16.1.248:2181"}
contexts := []string{"/configs/apps","/configs/users"}
cs = props.NewZookeeperCompositeConfigSource(contexts, urls, time.Second*3)

```

#### 用properties来配置： key/properties

value值为properties格式内容, 整体设计类似ini格式

```golang
root := "/config/kv/app1/dev"
var conn *zk.Conn
p := props.NewZookeeperIniConfigSource("zookeeper-kv", root, conn)

```


### consul 多层key/value形式


#### by consul key/value

```golang
例如：

config101/test/demo1/server/port=8080

获取的属性和值是：

server.port=8080

address := "127.0.0.1:8500"
root := "config101/test/demo1"
c := NewConsulKeyValueConfigSource("consul", address, root)
stringValue, err := cs.Get("prefix.key1")
stringDefaultValue := cs.GetDefault("prefix.key1", "default value")

```
#### 用properties来配置： key/properties

value值为properties格式内容, 整体设计类似ini格式

```golang
root := "config/app1/dev"
address := "127.0.0.1:8500"
p := props.NewConsulIniConfigSourceByName("consul-ini", address, root)
```

### 支持Unmarshal

在struct中规定命名为`_prefix `、类型为`string `、并且指定了`prefix`tag, 使用feild `_prefix `的`prefix`tag作为前缀，将struct feild名称转换后组合成完整的key，并从ConfigSource中获取数据并注入struct实例，feild类型只支持ConfigSource所支持的数据类型（string、int、float、bool、time.Duration）。

```golang

type ServerProperties struct {
	_prefix string `prefix:"http.server"`
	Port    int
	Timeout int `val:"1"`
	Enabled bool
	Foo     int `val:"1"`
	Time    time.Duration `val:"1s"`
	Float   float32 `val:"0.000001"`
}

func main() {

	p := props.NewMapProperties()
	p.Set("http.server.port", "8080")
	s := &ServerProperties{}
	p.Unmarshal(s)
	fmt.Println(s)

}


```

### 多种配置源组合使用

优先级以追加相反的顺序,最后添加优先级最高。

```golang

kv1 := []string{"go.app.key1", "value1", "value1-2"}
kv2 := []string{"go.app.key2", "value2", "value2-2"}

p1 := NewEmptyMapConfigSource("map1")
p1.Set(kv1[0], kv1[1])
p1.Set(kv2[0], kv2[1])
p2 := NewEmptyMapConfigSource("map2")
p2.Set(kv1[0], kv1[2])
p2.Set(kv2[0], kv2[2])
conf.Add(p1)
conf.Add(p2)

//value1==value1-2
value1, err := conf.Get(kv1[0])
fmt.Println(value1)
//value2=value2-2
value2, err := conf.Get(kv2[0])
fmt.Println(value2)

	
```

