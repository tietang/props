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
	Get(key string) (string, error)
	GetInt(key string) (int, error)
	GetBool(key string) (bool, error)
	GetFloat64(key string) (float64, error)
	Set(key, val string)
	SetAll(values map[string]string)
	Keys() []string
}

type PropertiesConfigSource struct {
	MapProperties
	name string
}

func NewPropertiesConfigSource(fileName string) *PropertiesConfigSource {
	return NewPropertiesConfigSourceByFile(fileName, fileName)
}

func NewPropertiesConfigSourceByFile(name, file string) *PropertiesConfigSource {
	p, err := ReadPropertyFile(file)
	var m map[string]string
	if err == nil {
		m = p.values
	} else {
		log.WithField("error", err.Error()).Info("read file: ")
	}
	s := &PropertiesConfigSource{}
	s.name = name
	s.values = m
	return s
}

func NewPropertiesConfigSourceByMap(name string, kv map[string]string) *PropertiesConfigSource {
	s := &PropertiesConfigSource{}
	s.name = name
	if kv == nil {
		s.values = make(map[string]string)
	} else {
		s.values = kv
	}
	return s
}

func (s *PropertiesConfigSource) Name() string {
	return s.name
}

//
type CompositeConfigSource struct {
	name          string
	ConfigSources []ConfigSource //Set
}

func NewCompositeConfigSource(mapSources []ConfigSource) *CompositeConfigSource {
	s := &CompositeConfigSource{
		ConfigSources: mapSources,
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
