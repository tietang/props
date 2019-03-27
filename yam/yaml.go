package yam

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"reflect"
)

type YamlProperties struct {
	kvs.MapProperties
}

func NewYamlProperties() *YamlProperties {
	p := &YamlProperties{}
	p.Values = make(map[string]string)
	return p
}

// Read creates a new property set and fills it with the contents of a file.
// See Load for the supported file format.
func ReadYaml(r io.Reader) (*YamlProperties, error) {
	p := NewYamlProperties()
	err := p.Load(r)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return p, nil
}

func ReadYamlFile(f string) (*YamlProperties, error) {

	file, err := os.Open(f)
	defer file.Close()

	if err != nil {
		d, _ := os.Getwd()
		log.WithField("error", err.Error()).Fatal("read file: ", d, "  ", f)
		return nil, err
	}
	return ReadYaml(file)
}

func (p *YamlProperties) Load(r io.Reader) error {

	data, err := ioutil.ReadAll(r)
	maps := make(map[string]interface{}, 0)
	err = yaml.Unmarshal([]byte(data), maps)
	if err != nil {
		log.Errorf("error: %v", err)
	}
	v := reflect.ValueOf(maps)
	p.kv(v, "")
	return nil
}

func (p *YamlProperties) kv(mapv reflect.Value, parentPath string) {
	iter := mapv.MapRange()
	for iter.Next() {
		k := iter.Key()
		val := iter.Value()
		valv := val.Elem()
		if valv.Kind() == reflect.Map {
			path := fmt.Sprintf("%s.%v", parentPath, k)
			p.kv(valv, path)
			continue
		} else {
			key := fmt.Sprintf("%s.%v", parentPath, k)[1:]
			value := fmt.Sprintf("%v", val)
			if value == "<nil>" {
				value = ""
			}
			p.Values[key] = value
		}
	}

}
