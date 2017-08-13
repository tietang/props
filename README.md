# props
config tools for golang, support properties file 、 zookeeper、consul


## 使用方法

## properties 文件读取

### 通过读取props.ReadPropertyFile

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

## ConfigSource 抽象方式读取

### Properties ConfigSource

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

### by zookeeper

```
var cs props.ConfigSource
urls := []string{"172.16.1.248:2181"}
contexts := []string{"/configs/apps","/configs/users"}
cs = props.NewZookeeperCompositeConfigSource(contexts, urls, time.Second*3)

```
### by consul key/value store

```
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

### 组合多个

```

	var pcs props.ConfigSource
	//通过文件名，文件名作为ConfigSource name
	pcs = props.NewPropertiesConfigSource("config.properties")



	//指定名称和文件名
	pcs2 := props.NewPropertiesConfigSourceByFile("config", "config.properties")
	//from zookeeper
	urls := []string{"172.16.1.248:2181"}
	contexts := []string{"/configs/apps", "/configs/users"}
	zccs := props.NewZookeeperCompositeConfigSource(contexts, urls, time.Second*3)
	configSources := []props.ConfigSource{pcs, pcs2, zccs, }
	ccs := props.NewDefaultCompositeConfigSource(configSources)

	
```