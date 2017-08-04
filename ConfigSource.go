package props

import (
	"errors"
	"strings"
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"regexp"
	log "github.com/sirupsen/logrus"
)

var reg = regexp.MustCompile("\\$\\{(.*)}")

type ConfigSource interface {
	Name() string
	//
	Get(key string) (string, error)
	GetInt(key string) (int, error)
	GetBool(key string) (bool, error)
	GetFloat64(key string) (float64, error)
	//
	GetDefault(key, defaultValue string) string
	GetIntDefault(key string, defaultValue int) int
	GetBoolDefault(key string, defaultValue bool) bool
	GetFloat64Default(key string, defaultValue float64) float64
	//
	Set(key, val string)
	SetAll(values map[string]string)
	Keys() []string
}

//
type CompositeConfigSource struct {
	name          string
	ConfigSources []ConfigSource //Set
}

func NewEmptyCompositeConfigSource() *CompositeConfigSource {
	s := &CompositeConfigSource{
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

func NewZookeeperCompositeConfigSource(contexts []string, connStr []string, timeout time.Duration) *CompositeConfigSource {

	conn, _, err := zk.Connect(connStr, timeout)
	if err != nil {
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

func (s *CompositeConfigSource) Name() string {
	return s.name
}
func (s *CompositeConfigSource) Add(ms ConfigSource) {
	for _, s := range s.ConfigSources {
		if ms.Name() == s.Name() {
			log.Info("exits ConfigSource: " + s.Name())
			return
		}
	}
	s.ConfigSources = append(s.ConfigSources, ms)

}

func (s *CompositeConfigSource) Get(key string) (string, error) {
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
	for _, s := range s.ConfigSources {
		v, err := s.Get(key)
		if err == nil {
			return v, nil
		}
	}
	return "", errors.New("not exists for key: " + key)
}
//
//func (s *CompositeConfigSource) evalValue(val string) string {
//	if strings.(val,)
//}

func (s *CompositeConfigSource) GetInt(key string) (int, error) {
	for _, s := range s.ConfigSources {
		v, err := s.GetInt(key)
		if err == nil {
			return v, nil
		}
	}
	return 0, errors.New("not exists for key: " + key)
}

func (s *CompositeConfigSource) GetBool(key string) (bool, error) {
	for _, s := range s.ConfigSources {
		v, err := s.GetBool(key)
		if err == nil {
			return v, nil
		}
	}
	return false, errors.New("not exists for key: " + key)
}

func (s *CompositeConfigSource) GetFloat64(key string) (float64, error) {
	for _, s := range s.ConfigSources {
		v, err := s.GetFloat64(key)
		if err == nil {
			return v, nil
		}
	}
	return 0.0, errors.New("not exists for key: " + key)
}

func (s *CompositeConfigSource) GetDefault(key string, defaultValue string) string {
	v, err := s.Get(key)
	if err != nil {
		return defaultValue
	}
	return v

}

func (s *CompositeConfigSource) GetIntDefault(key string, defaultValue int) int {
	v, err := s.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return v

}

func (s *CompositeConfigSource) GetBoolDefault(key string, defaultValue bool) bool {
	v, err := s.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return v
}

func (s *CompositeConfigSource) GetFloat64Default(key string, defaultValue float64) float64 {
	v, err := s.GetFloat64(key)
	if err != nil {
		return defaultValue
	}
	return v
}

func (s *CompositeConfigSource) Set(key, val string) {
	panic(errors.New("Unsupported operation"))
}

func (s *CompositeConfigSource) SetAll(values map[string]string) {
	panic(errors.New("Unsupported operation"))
}

func (s *CompositeConfigSource) Keys() []string {
	keys := make([]string, 0)
	for _, s := range s.ConfigSources {
		ks := s.Keys()
		for _, k := range ks {
			keys = append(keys, k)
		}

	}
	return keys
}

func (s *CompositeConfigSource) GetValue(key string) (string, error) {
	v, err := s.Get(key)
	if err == nil {
		if reg.MatchString(v) {
			return s.calculateValue(v)
		}
		return v, nil
	}
	return v, err
}

func (s *CompositeConfigSource) calculateValue(value string) (string, error) {
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
		v, err := s.Get(k)
		if err == nil {
			return v, nil
		}
	}

	return defaultValue, errors.New("not exists")
}
