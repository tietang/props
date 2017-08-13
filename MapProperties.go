package props

import (
	"strconv"
	"errors"
	"time"
	"strings"
	"reflect"
)

const (
	PREFIX_FIELD            = "_prefix"
	STRUCT_PREFIX_TAG       = "prefix"
	FIELD_DEFAULT_VALUE_TAG = "val"
)

type MapProperties struct {
	values map[string]string
}

func NewMapProperties() *MapProperties {
	p := &MapProperties{
		values: make(map[string]string),
	}
	return p
}

func NewMapPropertiesByMap(kv map[string]string) *MapProperties {
	p := &MapProperties{
		values: kv,
	}
	return p
}

func (p *MapProperties) Name() string {
	return "MapProperties"
}

//--get key/value

// Get retrieves the value of a property. If the property does not exist, an
// empty string will be returned.
func (p *MapProperties) Get(key string) (string, error) {
	if v, ok := p.values[key]; ok {
		return v, nil
	}
	return "", errors.New("not exists for key: " + key)
}

// GetDefault retrieves the value of a property. If the property does not
// exist, then the default value will be returned.
func (p *MapProperties) GetDefault(key, defVal string) string {
	if v, ok := p.values[key]; ok {
		return v
	}
	return defVal
}

func (p *MapProperties) GetInt(key string) (int, error) {
	v := p.values[key]

	if v, err := strconv.Atoi(v); err == nil {
		return v, nil
	}
	return 0, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetIntDefault(key string, defVal int) int {
	if v, ok := p.values[key]; ok {
		if v, err := strconv.Atoi(v); err == nil {
			return v
		}
	}
	return defVal
}
func (p *MapProperties) GetBool(key string) (bool, error) {
	v := p.values[key]

	if v, err := strconv.ParseBool(v); err == nil {
		return v, nil
	}
	return false, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetBoolDefault(key string, defVal bool) bool {
	if v, ok := p.values[key]; ok {
		if v, err := strconv.ParseBool(v); err == nil {
			return v
		}
	}
	return defVal
}

func (p *MapProperties) GetFloat64(key string) (float64, error) {
	v := p.values[key]

	if v, err := strconv.ParseFloat(v, 64); err == nil {
		return v, nil
	}
	return 0.0, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetFloat64Default(key string, defVal float64) float64 {
	if v, ok := p.values[key]; ok {
		if v, err := strconv.ParseFloat(v, 64); err == nil {
			return v
		}
	}
	return defVal
}

// 1ms 1mS 1MS 1Ms -> 1*time.Millisecond
//1s 1 1S -> 1*time.Second
//无单位默认为second
func (p *MapProperties) GetDuration(key string) (time.Duration, error) {
	v, err := p.Get(key)
	if err != nil {

		return time.Duration(0), err
	}
	return toDuration(v)
}

func (p *MapProperties) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	if v, ok := p.values[key]; ok {
		if v, err := toDuration(v); err == nil {
			return v
		}
	}
	return defaultValue
}

func toDuration(v string) (time.Duration, error) {

	v = strings.ToUpper(v)

	if strings.LastIndex(v, TIME_MS) > 0 {
		i, err := strconv.ParseInt(strings.TrimSuffix(v, TIME_MS), 10, 0)
		return time.Duration(i) * time.Millisecond, err
	} else {
		i, err := strconv.ParseInt(strings.TrimSuffix(v, TIME_S), 10, 0)
		return time.Duration(i) * time.Second, err
	}
}

// Names returns the keys for all properties in the set.
func (p *MapProperties) Keys() []string {
	keys := make([]string, 0, len(p.values))
	for k, _ := range p.values {
		keys = append(keys, k)
	}
	return keys
}

// Set adds or changes the value of a property.
func (p *MapProperties) Set(key, val string) {
	p.values[key] = val
}
func (p *MapProperties) SetAll(values map[string]string) {
	for k, v := range values {
		p.values[k] = v
	}

}

// Clear removes all key-value pairs.
func (p *MapProperties) Clear() {
	p.values = make(map[string]string)
}

func (p *MapProperties) Unmarshal(obj interface{}) error {
	return Unmarshal(p, obj)
}

func Unmarshal(p ConfigSource, obj interface{}) error {

	//t := reflect.TypeOf(obj)
	//num := t.NumField()
	//for i := 0; i < num; i++ {
	//	sf := t.Field(i)
	//	fmt.Println(sf.Name)
	//}
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()
	num := v.NumField()
	sf, ok := t.FieldByName(PREFIX_FIELD)
	prefix := ""
	if ok {
		//fmt.Println(err)
		//fmt.Println(sf.Name)
		prefix = sf.Tag.Get(STRUCT_PREFIX_TAG)
		//fmt.Println(prefix)
	}

	for i := 0; i < num; i++ {
		sf := t.Field(i)
		if sf.Name == PREFIX_FIELD {
			continue
		}
		keys := toKeys(sf.Name)
		key1 := strings.Join([]string{prefix, keys[0]}, ".")
		key2 := strings.Join([]string{prefix, keys[1]}, ".")
		//fmt.Println(sf.Name)
		defVal := sf.Tag.Get(FIELD_DEFAULT_VALUE_TAG)

		value := v.Field(i) // reflect.ValueOf(sf)
		if !value.IsValid() || !value.CanSet() {
			continue
		}
		//value := reflect.ValueOf(&value1).Elem()
		//fmt.Println("value: ", key1, key2, value.CanSet(), value.Type().Name(), value.Kind().String())
		switch value.Type().Name() {
		case "string":
			//fmt.Println(value)
			val1 := p.GetDefault(key1, defVal)
			val2 := p.GetDefault(key2, defVal)
			if val1 == "" {
				val1 = val2
			}
			value.SetString(val1)
			break
		case "int", "int32", "int64":
			val := getInt(p, key1, key2, defVal)
			//fmt.Println("setInt", val, value.CanSet(), reflect.ValueOf(&value).Elem().CanSet())
			value.SetInt(int64(val))
			break
		case "uint", "uint32", "uint64":
			val := getInt(p, key1, key2, defVal)
			//fmt.Println("-----", val)
			value.SetUint(uint64(val))
			break
		case "float32", "float64":

			defaultValue, err := strconv.ParseFloat(defVal, 64)
			if err != nil {
			}

			val1 := p.GetFloat64Default(key1, defaultValue)
			val2 := p.GetFloat64Default(key2, defaultValue)
			//fmt.Println("ff  ", sf.Name, "	", val1, val2)
			if val1 == 0 {
				val1 = val2
			}
			value.SetFloat(val1)
			break
		case "bool":
			defaultValue, err := strconv.ParseBool(defVal)
			if err != nil {
			}

			val1 := p.GetBoolDefault(key1, defaultValue)
			val2 := p.GetBoolDefault(key2, defaultValue)
			//fmt.Println("ff  ", sf.Name, "	", val1, val2)
			value.SetBool(val1 || val2)
			break
		case "Duration":
			defaultValue, err := toDuration(defVal)
			if err != nil {
			}

			val1 := p.GetDurationDefault(key1, defaultValue)
			val2 := p.GetDurationDefault(key2, defaultValue)
			//fmt.Println("duration:   ", defaultValue, val1, val2)

			if val2 <= 0 {
				val1 = val2
			}
			val1.Nanoseconds()
			value.SetInt(val1.Nanoseconds())
			break
		default:

		}
	}
	return nil
}
func getInt(p ConfigSource, key1, key2, defVal string) int {
	defaultValue, err := strconv.Atoi(defVal)
	if err != nil {
	}
	val1 := p.GetIntDefault(key1, defaultValue)
	val2 := p.GetIntDefault(key2, defaultValue)
	if val1 == 0 {
		val1 = val2
	}
	return val1

}

//
func toKeys(str string) [2]string {
	keys := [2]string{"", ""}
	keys[1] = strings.ToLower(str[0:1]) + str[1:]
	r := []rune(str)
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

	return keys
}
