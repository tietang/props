package props

import (
	log "github.com/sirupsen/logrus"
	"errors"
	"strings"
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"regexp"
	"github.com/valyala/fasttemplate"
	"io"
)

const (
	START_TAG     = "${"
	END_TAG       = "}"
	DEFAULT_VALUE = ""
)

var reg = regexp.MustCompile("\\$\\{(.*)}")

type ConfigSource interface {
	Name() string
	//
	Get(key string) (string, error)
	GetInt(key string) (int, error)
	GetDuration(key string) (time.Duration, error)
	GetBool(key string) (bool, error)
	GetFloat64(key string) (float64, error)
	//
	GetDefault(key, defaultValue string) string
	GetIntDefault(key string, defaultValue int) int
	GetDurationDefault(key string, defaultValue time.Duration) time.Duration

	GetBoolDefault(key string, defaultValue bool) bool
	GetFloat64Default(key string, defaultValue float64) float64
	//
	Set(key, val string)
	SetAll(values map[string]string)
	Keys() []string
	Unmarshal(t interface{}) error
}

//
type CompositeConfigSource struct {
	name          string
	ConfigSources []ConfigSource //Set
}

func NewEmptyCompositeConfigSource() *CompositeConfigSource {
	s := &CompositeConfigSource{
		ConfigSources: make([]ConfigSource, 0),
	}
	s.name = "CompositeConfigSource"

	return s
}

func NewDefaultCompositeConfigSource(configSources []ConfigSource) *CompositeConfigSource {
	s := &CompositeConfigSource{
		ConfigSources: configSources,
	}
	s.name = "CompositeConfigSource"

	return s
}

func NewCompositeConfigSource(name string, configSources []ConfigSource) *CompositeConfigSource {
	s := &CompositeConfigSource{
		ConfigSources: configSources,
	}
	s.name = "CompositeConfigSource"

	return s
}

func NewConsulKeyValueCompositeConfigSource(contexts []string, address string) *CompositeConfigSource {
	s := &CompositeConfigSource{}
	for _, context := range contexts {
		c := NewConsulKeyValueConfigSource("consul:"+context, address, context)
		s.Add(c)
	}
	s.name = "ConsulKevValue"
	return s
}

func NewZookeeperCompositeConfigSource(contexts []string, connStr []string, timeout time.Duration) *CompositeConfigSource {

	conn, _, err := zk.Connect(connStr, timeout)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return NewZookeeperCompositeConfigSourceByConn(contexts, conn)
}

func NewZookeeperCompositeConfigSourceByConn(contexts []string, conn *zk.Conn) *CompositeConfigSource {

	s := &CompositeConfigSource{}
	for _, context := range contexts {
		zkms := NewZookeeperConfigSource("zk:"+context, context, conn)
		s.Add(zkms)
	}
	s.name = "Zookeeper"

	return s
}

func (ccs *CompositeConfigSource) Name() string {
	return ccs.name
}
func (ccs *CompositeConfigSource) Size() int {
	return len(ccs.ConfigSources)
}
func (ccs *CompositeConfigSource) Add(ms ConfigSource) {
	for _, s := range ccs.ConfigSources {
		if ms.Name() == s.Name() {
			return
		}
	}
	ccs.ConfigSources = append(ccs.ConfigSources, ms)

}

func (ccs *CompositeConfigSource) Get(key string) (string, error) {
	//var value string;
	//var found bool;
	//s.ConfigSources.ForEach(func(mapSource interface{}, ok bool) {
	//    if ok {
	//        ms := mapSource.(ConfigSource)
	//        v, err := ms.Get(key)
	//        if err == nil {
	//            value = v
	//            found = ok
	//            return -1
	//        }
	//    }
	//})
	//
	//if found {
	//    return value, nil
	//}
	//
	for _, s := range ccs.ConfigSources {
		v, err := s.Get(key)
		if err == nil {
			return v, nil
		}
	}
	return "", errors.New("not exists for key: " + key)
}
func (ccs *CompositeConfigSource) GetInt(key string) (int, error) {
	for _, s := range ccs.ConfigSources {
		v, err := s.GetInt(key)
		if err == nil {
			return v, nil
		}
	}
	return 0, errors.New("not exists for key: " + key)
}
func (ccs *CompositeConfigSource) GetDuration(key string) (time.Duration, error) {
	for _, s := range ccs.ConfigSources {
		return s.GetDuration(key)
	}
	return time.Duration(0), errors.New("not exists for key: " + key)
}

func (ccs *CompositeConfigSource) GetBool(key string) (bool, error) {
	for _, s := range ccs.ConfigSources {
		v, err := s.GetBool(key)
		if err == nil {
			return v, nil
		}
	}
	return false, errors.New("not exists for key: " + key)
}

func (ccs *CompositeConfigSource) GetFloat64(key string) (float64, error) {
	for _, s := range ccs.ConfigSources {
		v, err := s.GetFloat64(key)
		if err == nil {
			return v, nil
		}
	}
	return 0.0, errors.New("not exists for key: " + key)
}

func (ccs *CompositeConfigSource) GetDefault(key string, defaultValue string) string {
	v, err := ccs.Get(key)
	if err != nil {
		return defaultValue
	}
	return v

}

func (ccs *CompositeConfigSource) GetIntDefault(key string, defaultValue int) int {
	v, err := ccs.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return v

}
func (ccs *CompositeConfigSource) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	v, err := ccs.GetDuration(key)
	if err != nil {
		return defaultValue
	}
	return v

}

func (ccs *CompositeConfigSource) GetBoolDefault(key string, defaultValue bool) bool {
	v, err := ccs.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return v
}

func (ccs *CompositeConfigSource) GetFloat64Default(key string, defaultValue float64) float64 {
	v, err := ccs.GetFloat64(key)
	if err != nil {
		return defaultValue
	}
	return v
}
func (ccs *CompositeConfigSource) Set(key, val string) {
	panic(errors.New("Unsupported operation"))
}

func (ccs *CompositeConfigSource) SetAll(values map[string]string) {
	panic(errors.New("Unsupported operation"))
}

func (ccs *CompositeConfigSource) Unmarshal(obj interface{}) error {
	return Unmarshal(ccs, obj)
}

func (ccs *CompositeConfigSource) Keys() []string {
	keys := make([]string, 0)
	for _, s := range ccs.ConfigSources {
		ks := s.Keys()
		for _, k := range ks {
			keys = append(keys, k)
		}

	}
	return keys
}

func (ccs *CompositeConfigSource) GetValue(key string) (string, error) {
	v, err := ccs.Get(key)
	if err == nil {
		if reg.MatchString(v) {
			return ccs.evalValue(v)
		}
		return v, nil
	}
	return v, err
}

func (ccs *CompositeConfigSource) evalValue(value string) (string, error) {
	if strings.Contains(value, START_TAG) {
		eval := fasttemplate.New(value, START_TAG, END_TAG)
		str := eval.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
			s, err := ccs.Get(tag)
			if err == nil {
				return w.Write([]byte(s))
			} else {
				return w.Write([]byte(""))
			}
		})
		return str, nil
	}
	return value, nil
}

func (ccs *CompositeConfigSource) calculateEvalValue(value string) (string, error) {

	sub := reg.FindStringSubmatch(value)
	if len(sub) == 0 {
		return value, nil
	}
	defaultValue := ""
	for _, k := range sub {

		keys := strings.Split(k, ":")
		if len(keys) > 1 {
			k = keys[0]
			defaultValue = keys[1]
		}
		v, err := ccs.Get(k)
		if err == nil {
			return v, nil
		}
	}

	return defaultValue, errors.New("not exists")
}
