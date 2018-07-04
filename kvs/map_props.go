package kvs

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
    Values map[string]string
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

//--get key/value
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
    v := p.Values[key]

    if v, err := strconv.ParseBool(v); err == nil {
        return v, nil
    }
    return false, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetBoolDefault(key string, defVal bool) bool {
    if v, ok := p.Values[key]; ok {
        if v, err := strconv.ParseBool(v); err == nil {
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
//1s 1 1S -> 1*time.Second
//无单位默认为second
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
    p.Values[key] = val
}
func (p *MapProperties) SetAll(values map[string]string) {
    for k, v := range values {
        p.Values[k] = v
    }

}

// Clear removes all key-value pairs.
func (p *MapProperties) Clear() {
    p.Values = make(map[string]string)
}

func (p *MapProperties) Unmarshal(obj interface{}) error {
    return Unmarshal(p, obj)
}

func Unmarshal(p ConfigSource, obj interface{}, parentKeys ...string) error {
    v := reflect.ValueOf(obj).Elem()
    return unmarshalInner(p, v, parentKeys...)
}
func unmarshalInner(p ConfigSource, v reflect.Value, parentKeys ...string) error {

    //t := reflect.TypeOf(obj)
    //num := t.NumField()
    //for i := 0; i < num; i++ {
    //	sf := t.Field(i)
    //	fmt.Println(sf.ConfName)
    //}

    t := v.Type()
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
        ks := toKeys(sf.Name)
        //fmt.Println(ks)
        keys := make([]string, 0)
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
                //fmt.Println(keys)

            }
        }

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
                val1 := p.GetDefault(key, defVal)
                if val1 != "" {
                    val = val1
                }
            }
            value.SetString(val)
            break
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

            if value.Type().Name() == "Duration" {
                defaultValue, err := ToDuration(defVal)
                if err != nil {
                }
                if value.Int() != 0 {
                    defaultValue = time.Nanosecond * time.Duration(value.Int())
                }
                val := defaultValue
                for _, key := range keys {
                    val1 := p.GetDurationDefault(key, defaultValue)
                    if val1 <= 0 {
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
                val1 := p.GetFloat64Default(key, defaultValue)
                if val1 != 0 {
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
                val1 := p.GetBoolDefault(key, defaultValue)
                if val1 {
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
            break
        case reflect.Struct:
            //fmt.Println("---")
            unmarshalInner(p, value, keys...)
            break
        default:

        }
    }
    return nil
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
        val1 := p.GetIntDefault(key, int(defaultValue))
        if val1 != 0 {
            val = int64(val1)
        }
    }
    return int(val)

}

//
func toKeys(str string) [2]string {
    keys := [2]string{"", ""}
    keys[1] = strings.ToLower(str[0:1]) + str[1:]
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
