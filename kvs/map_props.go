package kvs

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	PREFIX_FIELD            = "_prefix"
	STRUCT_PREFIX_TAG       = "prefix"
	FIELD_CONFIG_NAME_TAG   = "props"
	FIELD_DEFAULT_VALUE_TAG = "val"
)

var fieldConfigNameTags = [...]string{"props", "yaml", "yam", "json", "ini", "toml"}
var _ ConfigSource = new(MapProperties)

type UnmarshalListener struct {
	Prefixes []string
	Obj      any
}
type MapProperties struct {
	Values             map[string]string
	OnChanges          map[string][]func(k, v string)
	unmarshalListeners []*UnmarshalListener
}

func NewMapProperties() *MapProperties {
	p := &MapProperties{
		Values: make(map[string]string),
	}
	return p
}

func NewMapPropertiesByMap(kv map[string]string) *MapProperties {
	p := &MapProperties{
		Values: kv,
	}
	return p
}

func (p *MapProperties) Name() string {
	return "MapProperties"
}
func (p *MapProperties) AddChangeListener(key string, listener func(k, v string)) {
	if p.OnChanges == nil {
		p.OnChanges = make(map[string][]func(k string, v string))
	}
	listeners, found := p.OnChanges[key]
	if !found {
		listeners = make([]func(k, v string), 0, 8)
	}
	listeners = append(listeners, listener)
	p.OnChanges[key] = listeners
}

func (p *MapProperties) addChangeUnmarshalListener(obj interface{}, prefixes ...string) {
	l := &UnmarshalListener{
		Prefixes: prefixes,
		Obj:      obj,
	}
	p.unmarshalListeners = append(p.unmarshalListeners, l)
}

func (p *MapProperties) unmarshalAllListeners() {
	go time.AfterFunc(3*time.Second, func() {
		for _, listener := range p.unmarshalListeners {
			p.Unmarshal(listener.Obj, listener.Prefixes...)
		}
	})
}

// --get key/value
func (p *MapProperties) KeyValue(key string) *KeyValue {
	v := p.GetDefault(key, "")
	kv := NewKeyValue(key, v)
	return kv
}
func (p *MapProperties) Strings(key string) []string {
	return p.KeyValue(key).Strings()
}
func (p *MapProperties) Ints(key string) []int {
	return p.KeyValue(key).Ints()
}
func (p *MapProperties) Float64s(key string) []float64 {
	return p.KeyValue(key).Float64s()
}
func (p *MapProperties) Durations(key string) []time.Duration {
	return p.KeyValue(key).Durations()
}

// Get retrieves the value of a property. If the property does not exist, an
// empty string will be returned.
func (p *MapProperties) Get(key string) (string, error) {
	if v, ok := p.Values[key]; ok {
		return v, nil
	}
	return "", errors.New("not exists for key: " + key)
}

// GetDefault retrieves the value of a property. If the property does not
// exist, then the default value will be returned.
func (p *MapProperties) GetDefault(key, defVal string) string {
	if v, ok := p.Values[key]; ok {
		return v
	}
	return defVal
}

func (p *MapProperties) GetInt(key string) (int, error) {
	v := p.Values[key]

	if v, err := strconv.Atoi(v); err == nil {
		return v, nil
	}
	return 0, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetIntDefault(key string, defVal int) int {
	if v, ok := p.Values[key]; ok {
		if v, err := strconv.Atoi(v); err == nil {
			return v
		}
	}
	return defVal
}
func (p *MapProperties) GetBool(key string) (bool, error) {
	v, ok := p.Values[key]
	if ok {
		return ParseBool(strings.TrimSpace(v))
	}
	return false, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetBoolDefault(key string, defVal bool) bool {
	if v, ok := p.Values[key]; ok {
		if v, err := ParseBool(strings.TrimSpace(v)); err == nil {
			return v
		}
	}
	return defVal
}

func (p *MapProperties) GetFloat64(key string) (float64, error) {
	v := p.Values[key]

	if v, err := strconv.ParseFloat(v, 64); err == nil {
		return v, nil
	}
	return 0.0, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetFloat64Default(key string, defVal float64) float64 {
	if v, ok := p.Values[key]; ok {
		if v, err := strconv.ParseFloat(v, 64); err == nil {
			return v
		}
	}
	return defVal
}

// 1ms 1mS 1MS 1Ms -> 1*time.Millisecond
// 1s 1 1S -> 1*time.Second
// 无单位默认为second
func (p *MapProperties) GetDuration(key string) (time.Duration, error) {
	v, err := p.Get(key)
	if err != nil {

		return time.Duration(0), err
	}
	return ToDuration(v)
}

func (p *MapProperties) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	if v, ok := p.Values[key]; ok {
		if v, err := ToDuration(v); err == nil {
			return v
		}
	}
	return defaultValue
}

func (p *MapProperties) GetTime(key string) (time.Time, error) {
	v, err := p.Get(key)
	if err != nil {
		return cast.ToTime(0), err
	}
	if x, ok := IsNumInt(v); ok {
		return cast.ToTimeE(x)
	}
	return cast.ToTimeE(v)
}

func (p *MapProperties) GetTimeDefault(key string, defaultValue time.Time) time.Time {
	v, err := p.GetTime(key)
	if err == nil {
		return v
	} else {
		return defaultValue
	}

}

// Names returns the keys for all Properties in the set.
func (p *MapProperties) Keys() []string {
	keys := make([]string, 0, len(p.Values))
	for k, _ := range p.Values {
		keys = append(keys, k)
	}
	return keys
}

// Set adds or changes the value of a property.
func (p *MapProperties) Set(key, val string) {
	CheckAndPublishChange(p.Values, key, val, p.OnChanges)
	p.Values[key] = val
	p.unmarshalAllListeners()
}

func (p *MapProperties) SetAll(values map[string]string) {
	for k, v := range values {
		CheckAndPublishChange(p.Values, k, v, p.OnChanges)
		p.Values[k] = v
	}
	p.unmarshalAllListeners()
}

// Clear removes all key-value pairs.
func (p *MapProperties) Clear() {
	p.Values = make(map[string]string)
}

func (p *MapProperties) Unmarshal(obj interface{}, prefixes ...string) error {
	p.addChangeUnmarshalListener(obj, prefixes...)
	return Unmarshal(p, obj, prefixes...)
}

func Unmarshal(p ConfigSource, obj interface{}, parentKeys ...string) error {

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Struct {
		v = v.Elem()
	}
	if v.Kind() == reflect.Map {
		return unmarshalInnerMap(p, v, parentKeys...)
	}
	return unmarshalInner(p, v, parentKeys...)
}
func unmarshalInnerMap(p ConfigSource, v reflect.Value, parentKeys ...string) (err error) {
	keys := p.Keys()
	t := v.Type()
	typ := t.Elem()
	//fmt.Println(typ, typ.Kind())

	for _, key := range keys {
		for _, pkey := range parentKeys {
			if strings.HasPrefix(key, pkey) {
				k := key[len(pkey)+1:]
				idx := strings.Index(k, ".")
				k = k[:idx]
				pk := pkey + "." + k
				var mvalue reflect.Value

				if typ.Kind() == reflect.Ptr {
					mvalue = reflect.New(typ.Elem())
				}
				if typ.Kind() == reflect.Struct {
					mvalue = reflect.New(typ)
				}
				err := unmarshalInner(p, mvalue.Elem(), pk)
				if err != nil {
					log.Error(err)
				}
				//fmt.Println(mvalue.Elem())
				if typ.Kind() == reflect.Ptr {
					v.SetMapIndex(reflect.ValueOf(k), mvalue)
				}
				if typ.Kind() == reflect.Struct {
					v.SetMapIndex(reflect.ValueOf(k), mvalue.Elem())
				}

			}
		}
	}

	return nil
}

func unmarshalInner(p ConfigSource, v reflect.Value, parentKeys ...string) (err error) {

	//t := reflect.TypeOf(obj)
	//num := t.NumField()
	//for i := 0; i < num; i++ {
	//	sf := t.Field(i)
	//	fmt.Println(sf.ConfName)
	//}

	t := v.Type()
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	num := v.NumField()
	sf, ok := t.FieldByName(PREFIX_FIELD)
	prefix := ""

	if ok && (parentKeys == nil || len(parentKeys) == 0) {
		//fmt.Println(err)
		//fmt.Println(sf.ConfName)
		prefix = sf.Tag.Get(STRUCT_PREFIX_TAG)
		//fmt.Println("prefix: ", prefix)
	}

	prefix = strings.TrimSpace(prefix)
	for i := 0; i < num; i++ {
		sf := t.Field(i)
		if sf.Name == PREFIX_FIELD {
			continue
		}
		//configKey := sf.Tag.Get(FIELD_CONFIG_NAME_TAG)
		//ks := toKeys1(sf.Name, configKey)
		ks := toKeys(sf)
		//fmt.Println(ks)
		keys := make([]string, 0)
		if sf.Anonymous {
			keys = parentKeys
		} else {
			for _, k := range ks {
				if k == "" {
					continue
				}
				if parentKeys != nil && len(parentKeys) > 0 {
					for _, pk := range parentKeys {
						keys = append(keys, strings.Join([]string{pk, k}, "."))
					}
				} else {
					if prefix != "" {
						keys = append(keys, strings.Join([]string{prefix, k}, "."))
					} else {
						keys = append(keys, k)
					}

				}
			}
		}
		//fmt.Println(keys)
		//key1 := strings.Join([]string{prefix, keys[0]}, ".")
		//key2 := strings.Join([]string{prefix, keys[1]}, ".")
		//fmt.Println(sf.ConfName)
		defVal := sf.Tag.Get(FIELD_DEFAULT_VALUE_TAG)

		value := v.Field(i) // reflect.ValueOf(sf)
		if !value.IsValid() || !value.CanSet() {
			continue
		}
		//value := reflect.ValueOf(&value1).Elem()
		//fmt.Println("value: ", keys, value.CanSet(), value.Type().ConfName(), value.Kind().String())
		//switch value.Type().ConfName() {
		//fmt.Println(value.Type().Kind() == value.Kind())
		switch value.Kind() {
		case reflect.String:
			//case "string":
			if value.String() != "" {
				defVal = value.String()
			}
			val := defVal

			for _, key := range keys {
				//val1 := p.GetDefault(key, defVal)
				//if val1 != "" {
				//	val = val1
				//}
				val1, err := p.Get(key)
				//fmt.Printf("key: %s, val1: %s, err: %v\n", key, val1, err)
				if err == nil && val1 != "" {
					val = val1
				}
			}

			value.SetString(val)
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

			if value.Type().Name() == "Duration" {
				var defaultValue time.Duration
				if defVal == "" {
					defaultValue = time.Duration(0)
				} else {
					defaultValue, err = ToDuration(defVal)
					if err != nil {
						log.Warn(err)
					}
				}

				if value.Int() != 0 {
					defaultValue = time.Nanosecond * time.Duration(value.Int())
				}
				val := defaultValue
				for _, key := range keys {
					//val1 := p.GetDurationDefault(key, defaultValue)
					val1, err := p.GetDuration(key)
					//fmt.Printf("key: %s, val1: %d, err: %v\n", key, val1, err)
					if err == nil && val1 >= 0 {
						val = val1
					}

				}

				value.SetInt(val.Nanoseconds())

			} else {
				//case "int", "int32", "int64":
				val := getInt(p, keys, value.Int(), defVal)

				//fmt.Println("setInt", val, value.CanSet(), reflect.ValueOf(&value).Elem().CanSet())
				value.SetInt(int64(val))
			}

			break
			//case "uint", "uint32", "uint64":
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val := getInt(p, keys, int64(value.Uint()), defVal)
			//fmt.Println("-----", val)
			value.SetUint(uint64(val))
			break
			//case "float32", "float64":
		case reflect.Float32, reflect.Float64:

			defaultValue, err := strconv.ParseFloat(defVal, 64)
			if err != nil {
			}
			if value.Float() != 0 {
				defaultValue = value.Float()
			}

			val := defaultValue

			for _, key := range keys {
				//val1 := p.GetFloat64Default(key, defaultValue)
				//if val1 != 0 {
				//	val = val1
				//}
				val1, err := p.GetFloat64(key)
				//fmt.Printf("key: %s, val1: %f, err: %v\n", key, val1, err)
				if err == nil {
					val = val1
				}
			}
			value.SetFloat(val)
			break
			//case "bool":
		case reflect.Bool:
			defaultValue, err := strconv.ParseBool(defVal)
			if err != nil {
			}
			if value.Bool() {
				defaultValue = value.Bool()
			}

			val := defaultValue
			for _, key := range keys {
				//val1 := p.GetBoolDefault(key, defaultValue)
				//if val1 {
				//	val = val1
				//}
				val1, err := p.GetBool(key)
				//fmt.Printf("key: %s, val1: %v, err: %v\n", key, val1, err)
				if err == nil {
					val = val1
				}
			}
			value.SetBool(val)

			break
		case reflect.Map:

			t := value.Type()
			typ := t.Elem()
			if value.IsNil() || !value.IsValid() {
				value.Set(reflect.MakeMap(value.Type()))
			}
			for _, key := range p.Keys() {
				for _, k := range keys {
					if strings.HasPrefix(key, k) {
						mk := strings.TrimPrefix(key, k+".")
						mv := p.GetDefault(key, defVal)
						kv := NewKeyValue(mk, mv)
						v, err := marshalSimple(kv, typ)
						if err == nil {

							value.SetMapIndex(reflect.ValueOf(mk), reflect.ValueOf(v))
						}
					}
				}
			}
			break
		case reflect.Struct:
			//fmt.Println("---")
			err = unmarshalInner(p, value, keys...)
			break
		default:

		}
	}
	return err
}

func marshalSimple(kv *KeyValue, typ reflect.Type) (interface{}, error) {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ.Name() == "Duration" {
			return kv.Duration()
		} else {
			return kv.Int64()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return kv.Uint64()
	case reflect.Float32, reflect.Float64:
		return kv.Float64()
	case reflect.String:
		return kv.String(), nil
	case reflect.Bool:
		return kv.Bool()
	}
	return "", nil
}

func getInt(p ConfigSource, keys []string, originValue int64, defVal string) int {
	defaultValue := originValue
	defV, err := strconv.Atoi(defVal)
	if err != nil {
	}
	if defaultValue == 0 {
		defaultValue = int64(defV)
	}

	val := defaultValue

	for _, key := range keys {
		//val1 := p.GetIntDefault(key, int(defaultValue))
		//if val1 != 0 {
		//	val = int64(val1)
		//}
		val1, err := p.GetInt(key)
		//fmt.Printf("key: %s, val1: %d, err: %v\n", key, val1, err)
		if err == nil {
			val = int64(val1)
		}
	}

	return int(val)

}

func toKeys1(str string, configKey string) [3]string {
	keys := [3]string{"", ""}
	keys[1] = strings.ToLower(str[0:1]) + str[1:]
	keys[2] = configKey
	r := []rune(str)
	//     if strings.Index(str, "-") >= 0 {
	for i := 0; i < len(str); i++ {
		if i == 0 {
			keys[0] += strings.ToLower(string(r[i])) // + string(vv[i+1])
		} else {
			if r[i] >= 65 && r[i] < 91 {
				keys[0] += "-"
			}
			keys[0] += strings.ToLower(string(r[i]))
		}
	}
	//     }
	return keys
}

// 比如：`ServerPort` 可以映射为 `server_port` 、`ServerPort`、 `serverPort` 、 `server-port`、 `SERVER_PORT` 、 `SERVER-PORT` 等
func toKeys(field reflect.StructField) []string {
	keys := make([]string, 0)
	str := field.Name

	//if len(keys) >= 3 {
	//	return keys
	//}
	r := []rune(str)
	//     if strings.Index(str, "-") >= 0 {
	keys0 := [4]string{}

	for i := 0; i < len(str); i++ {
		if i == 0 {
			keys0[0] += strings.ToUpper(string(r[i])) // + string(vv[i+1])
			keys0[1] += strings.ToUpper(string(r[i]))
			keys0[2] += strings.ToLower(string(r[i]))
			keys0[3] += strings.ToLower(string(r[i]))
		} else {
			if r[i] >= 65 && r[i] < 91 {
				keys0[2] += "_"
				keys0[3] += "-"
				keys0[0] += "_"
				keys0[1] += "-"
			}
			keys0[0] += strings.ToUpper(string(r[i]))
			keys0[1] += strings.ToUpper(string(r[i]))
			keys0[2] += strings.ToLower(string(r[i]))
			keys0[3] += strings.ToLower(string(r[i]))
		}
	}
	keys = append(keys, keys0[:]...)
	keys = append(keys, str)
	keys = append(keys, strings.ToLower(str[0:1])+str[1:])
	for _, tag := range fieldConfigNameTags {
		configKey := field.Tag.Get(tag)
		if configKey != "" {
			keys = append(keys, configKey)
		}
	}

	//fmt.Println("keys: ", keys)
	return keys
}
