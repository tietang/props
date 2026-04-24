package tom

import (
	"fmt"

	"io"
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
)

type TomlProperties struct {
	kvs.MapProperties
}

func NewTomlProperties() *TomlProperties {
	p := &TomlProperties{}
	p.Values = make(map[string]string)
	return p
}

// Read creates a new property set and fills it with the contents of a file.
// See Load for the supported file format.
func ReadToml(r io.Reader) (*TomlProperties, error) {
	p := NewTomlProperties()
	err := p.Load(r)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return p, nil
}

func ReadTomlFile(f string) (*TomlProperties, error) {

	file, err := os.Open(f)
	defer file.Close()

	if err != nil {
		d, _ := os.Getwd()
		log.WithField("error", err.Error()).Fatal("read file: ", d, "  ", f)
		return nil, err
	}
	return ReadToml(file)
}

type Config map[string]interface{}

func (p *TomlProperties) Load(r io.Reader) error {

	data, err := io.ReadAll(r)
	//maps := make(map[string]interface{}, 0)
	////var maps interface{} = make(map[string]interface{}, 0)
	//err = toml.Unmarshal(data, maps)

	//var maps map[string]interface{}
	var maps Config
	err = toml.Unmarshal(data, &maps)
	if err != nil {
		log.Errorf("error: %v", err)
	}
	v := reflect.ValueOf(maps)
	p.kv(v, "")
	return nil
}

func (p *TomlProperties) kv(v reflect.Value, parentPath string) {

	if v.Kind() == reflect.Interface {
		if v.Elem().Kind() == reflect.Map {
			//p.kv(v.Elem(), parentPath)
			log.Warn("nested map not supported:", parentPath)
			return
		}
		if v.Elem().Kind() == reflect.Slice {
			log.Warn("nested slice not supported:", parentPath)
			return
		}
		////fmt.Printf("---- %s=%v\n", parentPath, v)
		//key := parentPath[1:]
		//value := fmt.Sprintf("%v", v)
		//if value == "<nil>" {
		//	value = ""
		//}
		//p.Values[key] = value
	}

	if v.Kind() == reflect.Map {
		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			val := iter.Value()
			valv := val.Elem()
			if valv.Kind() == reflect.Map {
				path := fmt.Sprintf("%s.%v", parentPath, k)
				p.kv(valv, path)
				continue
			} else if valv.Kind() == reflect.Slice {
				key := fmt.Sprintf("%s.%v", parentPath, k)
				var vals string
				for i := 0; i < valv.Len(); i++ {
					v := valv.Index(i)
					if v.Kind() == reflect.Slice {
						log.Warn("nested slice not supported:", key)
						continue
					}
					if v.Kind() == reflect.Map {
						log.Warn("A slice nested within a map is not supported:", key)
						continue
					}
					if v.Kind() == reflect.Interface && v.Elem().Kind() == reflect.Map {
						log.Warn("nested slice not supported:", key)
						continue
					}
					if v.Kind() == reflect.Interface && v.Elem().Kind() == reflect.Slice {
						log.Warn("A slice nested within a map is not supported:", key)
						continue
					}
					vals += fmt.Sprintf("%s%v", kvs.DEFAULT_DELIMS, v)
				}
				if vals == "" {
					continue
				}
				p.Values[key[1:]] = vals[1:]
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

}

func (p *TomlProperties) kv2(v reflect.Value, parentPath string) {

	if v.Kind() == reflect.Interface {
		if v.Elem().Kind() == reflect.Map {
			p.kv(v.Elem(), parentPath)
			return
		}
		if v.Elem().Kind() == reflect.Slice {
			for i := 0; i < v.Elem().Len(); i++ {
				path := fmt.Sprintf("%s[%d]", parentPath, i)
				v := v.Elem().Index(i)
				p.kv(v, path)
			}
			return
		}
		//fmt.Printf("---- %s=%v\n", parentPath, v)
		key := parentPath[1:]
		value := fmt.Sprintf("%v", v)
		if value == "<nil>" {
			value = ""
		}
		p.Values[key] = value
	}

	if v.Kind() == reflect.Map {

		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			val := iter.Value()
			valv := val.Elem()
			if valv.Kind() == reflect.Map {
				path := fmt.Sprintf("%s.%v", parentPath, k)
				p.kv(valv, path)
				continue
			}
			if valv.Kind() == reflect.Slice {
				for i := 0; i < valv.Len(); i++ {
					path := fmt.Sprintf("%s.%v[%d]", parentPath, k, i)
					v := valv.Index(i)
					p.kv(v, path)
				}
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

}
