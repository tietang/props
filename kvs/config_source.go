package kvs

import (
	"bufio"
	"github.com/tietang/props/v3/kvs/helper"
	"path/filepath"
	"strings"
	"time"
)

type ContentType string

//配置为ContentIniProps模式，
// key作为section，value为props格式内容，类似ini文件格式
// key为实际key的prefix，会添加到前面
// 比如 root=configs/dev/app
// consul
// 		full key=configs/dev/app/mysql
// 		value=(x1=0 \n x2=1)
// 实际key/value为： mysql.x1=0 mysql.x2=1
//ContentProps,ContentYamlContentIni 模式时，
// 其 key无实际配置意义，只作为配置分组标识，
// 值为对应的内容格式类型，读取时会将对应的内容转换这种类型
// 可以通过key后缀来标识格式类型，默认按照properties来读取

const (
	ContentProperties   ContentType = "properties"
	ContentProps        ContentType = "props" //properties 别名
	ContentYaml         ContentType = "yaml"
	ContentYam          ContentType = "yam" //yaml 别名
	ContentYml          ContentType = "yml" //yaml 别名
	ContentIni          ContentType = "ini"
	ContentIniProps     ContentType = "ini_props"
	ContentKV           ContentType = "kv"
	ContentAuto         ContentType = "auto"
	TextContentType     ContentType = "text"
	TxtContentType      ContentType = "txt"
	KeyValueContentType ContentType = "kv"
	ContentUnknown      ContentType = "unknown"
	DefaultContentType  ContentType = ContentProps
)

var supportedPrefixChars = []string{"#@", ";@", "//@", "@"}

func ReadContentType(content string) ContentType {
	r := bufio.NewReader(strings.NewReader(content))
	i := 0
	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
			contentType, done := ReadLineContentType(string(line))
			if done {
				return contentType
			}
			i++
		}
		if err != nil || i > 0 {
			break
		}
	}
	return TextContentType
}

func ReadLineContentType(line string) (ContentType, bool) {
	str := strings.TrimSpace(line)
	for _, c := range supportedPrefixChars {
		if strings.HasPrefix(str, c) {
			ctyp := strings.TrimSpace(str[len(c):])
			if ctyp == "" {
				return TextContentType, true
			} else {
				return ContentType(str[len(c):]), true
			}
		}
	}
	return "", false
}

func GetContentTypeByName(name string) (ctype ContentType) {
	ctype = ContentUnknown
	idx := strings.LastIndex(name, ".")
	if !(idx == -1 || idx == len(name)-1) {
		ctype = ContentType(name[idx+1:])
		return
	}
	base := filepath.Base(name)
	names := helper.SplitAny(base, "._-")
	if len(names) > 1 {
		typName := names[len(names)-1]
		ctype = ContentType(typName)
	}

	return
}

type ConfigSource interface {
	Name() string
	//
	AddChangeListener(key string, listener func(k, v string))
	//
	KeyValue(key string) *KeyValue
	Strings(key string) []string
	Ints(key string) []int
	Float64s(key string) []float64
	Durations(key string) []time.Duration
	//
	Get(key string) (string, error)
	GetDefault(key, defaultValue string) string

	//
	GetInt(key string) (int, error)
	GetIntDefault(key string, defaultValue int) int
	//
	GetDuration(key string) (time.Duration, error)
	GetDurationDefault(key string, defaultValue time.Duration) time.Duration
	//
	GetTime(key string) (time.Time, error)
	GetTimeDefault(key string, defaultValue time.Time) time.Time
	//
	GetBool(key string) (bool, error)
	GetBoolDefault(key string, defaultValue bool) bool
	//
	GetFloat64(key string) (float64, error)
	GetFloat64Default(key string, defaultValue float64) float64
	//
	Set(key, val string)
	SetAll(values map[string]string)
	Keys() []string
	//KeysFilter(filter string) []string
	//t必须为指针型
	Unmarshal(obj interface{}, prefixes ...string) error
}
