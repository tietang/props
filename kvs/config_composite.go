package kvs

import (
    "time"
    "errors"
    "strings"
    "github.com/valyala/fasttemplate"
    "io"
)

//
type CompositeConfigSource struct {
    ConfName      string
    ConfigSources []ConfigSource //Set
    Properties    *PropertiesConfigSource
    IsEval        bool
    StartTag      string
    EndTag        string
}

func NewEmptyCompositeConfigSource() *CompositeConfigSource {
    return NewCompositeConfigSource("CompositeConfigSource", true)
}
func NewEmptyNoSystemEnvCompositeConfigSource() *CompositeConfigSource {
    return NewCompositeConfigSource("CompositeConfigSource-NoSystemEnv", false)
}

func NewDefaultCompositeConfigSource(configSources ...ConfigSource) *CompositeConfigSource {
    return NewCompositeConfigSource("CompositeConfigSource", true, configSources...)
}

func NewCompositeConfigSource(name string, isAppendSystemEnv bool, configSources ...ConfigSource) *CompositeConfigSource {
    s := &CompositeConfigSource{
        ConfigSources: make([]ConfigSource, 0),
        ConfName:      name,
        Properties:    NewEmptyMapConfigSource("default_properties"),
        StartTag:      __START_TAG,
        EndTag:        __END_TAG,
    }
    if name == "" {
        s.ConfName = "CompositeConfigSource"
    }
    if isAppendSystemEnv {
        s.ConfigSources = append(s.ConfigSources, s.Properties, newEnvConfigSource())
    } else {
        s.ConfigSources = append(s.ConfigSources, s.Properties)
    }

    for _, cs := range configSources {
        s.ConfigSources = append(s.ConfigSources, cs)
    }
    return s
}

func (ccs *CompositeConfigSource) Name() string {
    return ccs.ConfName
}
func (ccs *CompositeConfigSource) Size() int {
    return len(ccs.ConfigSources)
}

func (ccs *CompositeConfigSource) Add(css ...ConfigSource) {
    for _, conf := range css {
        for i := len(ccs.ConfigSources) - 1; i >= 0; i-- {
            s := ccs.ConfigSources[i]
            if conf.Name() == s.Name() {
                return
            }
        }
        ccs.ConfigSources = append(ccs.ConfigSources, conf)
    }
}

func (ccs *CompositeConfigSource) AddAll(css []ConfigSource) {
    ccs.Add(css...)
}

func (ccs *CompositeConfigSource) KeyValue(key string) *KeyValue {
    //v := ccs.GetDefault(key, "")
    //kv := NewKeyValue(key, v)
    //return kv

    val := ""
    hasExists := false
    for i := len(ccs.ConfigSources) - 1; i >= 0; i-- {
        s := ccs.ConfigSources[i]
        v, err := s.Get(key)
        if err == nil {
            val = v
            hasExists = true
            break
        }
    }

    if __reg.MatchString(val) {
        v, err := ccs.evalValue(val)
        kv := NewKeyValue(key, v)
        kv.err = err
        return kv
    }
    if hasExists {
        kv := NewKeyValue(key, val)
        kv.err = nil
        return kv
    } else {
        kv := NewKeyValue(key, val)
        kv.err = errors.New("not exists for key: " + key)
        return kv
    }
}
func (ccs *CompositeConfigSource) Strings(key string) []string {
    return ccs.KeyValue(key).Strings()
}
func (ccs *CompositeConfigSource) Ints(key string) []int {
    return ccs.KeyValue(key).Ints()
}
func (ccs *CompositeConfigSource) Float64s(key string) []float64 {
    return ccs.KeyValue(key).Float64s()
}

func (ccs *CompositeConfigSource) Durations(key string) []time.Duration {
    return ccs.KeyValue(key).Durations()
}

func (ccs *CompositeConfigSource) Get(key string) (string, error) {
    return ccs.GetValue(key)
}

func (ccs *CompositeConfigSource) GetDefault(key string, defaultValue string) string {
    v, err := ccs.GetValue(key)
    if err != nil {
        return defaultValue
    }
    return v
}
func (ccs *CompositeConfigSource) GetInt(key string) (int, error) {
    val, err := ccs.GetValue(key)
    if err == nil {
        return NewKeyValue(key, val).Int()
    } else {
        return 0, err
    }

}
func (ccs *CompositeConfigSource) GetDuration(key string) (time.Duration, error) {
    val, err := ccs.GetValue(key)
    if err == nil {
        return NewKeyValue(key, val).Duration()
    } else {
        return time.Duration(0), err
    }
}

func (ccs *CompositeConfigSource) GetBool(key string) (bool, error) {

    val, err := ccs.GetValue(key)
    if err == nil {
        return NewKeyValue(key, val).Bool()
    } else {
        return false, err
    }
}

func (ccs *CompositeConfigSource) GetFloat64(key string) (float64, error) {
    val, err := ccs.GetValue(key)
    if err == nil {
        return NewKeyValue(key, val).Float64()
    } else {
        return 0.0, err
    }
}

func (ccs *CompositeConfigSource) GetIntDefault(key string, defaultValue int) int {

    val, err := ccs.GetValue(key)
    if err == nil {
        v, err := NewKeyValue(key, val).Int()
        if err == nil {
            return v
        } else {
            return defaultValue
        }
    } else {
        return defaultValue
    }

}
func (ccs *CompositeConfigSource) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {

    val, err := ccs.GetValue(key)
    if err == nil {
        v, err := NewKeyValue(key, val).Duration()
        if err == nil {
            return v
        } else {
            return defaultValue
        }
    } else {
        return defaultValue
    }

}

func (ccs *CompositeConfigSource) GetBoolDefault(key string, defaultValue bool) bool {

    val, err := ccs.GetValue(key)
    if err == nil {
        v, err := NewKeyValue(key, val).Bool()
        if err == nil {
            return v
        } else {
            return defaultValue
        }
    } else {
        return defaultValue
    }
}

func (ccs *CompositeConfigSource) GetFloat64Default(key string, defaultValue float64) float64 {

    val, err := ccs.GetValue(key)
    if err == nil {
        v, err := NewKeyValue(key, val).Float64()
        if err == nil {
            return v
        } else {
            return defaultValue
        }
    } else {
        return defaultValue
    }
}
func (ccs *CompositeConfigSource) Set(key, val string) {
    //panic(errors.New("Unsupported operation"))
    ccs.Properties.Set(key, val)
}

func (ccs *CompositeConfigSource) SetAll(values map[string]string) {
    //panic(errors.New("Unsupported operation"))
    ccs.Properties.SetAll(values)
}

func (ccs *CompositeConfigSource) Unmarshal(obj interface{}) error {
    return Unmarshal(ccs, obj)
}

func (ccs *CompositeConfigSource) Keys() []string {

    set := NewSet()
    for i := len(ccs.ConfigSources) - 1; i >= 0; i-- {
        s := ccs.ConfigSources[i]
        ks := s.Keys()
        for _, k := range ks {
            set.Add(k)
        }
    }
    keys := make([]string, 0)
    set.ForEach(func(i interface{}, i2 bool) int {
        keys = append(keys, i.(string))
        return 1
    })
    return keys
}

func (ccs *CompositeConfigSource) GetValue(key string) (string, error) {
    //val := ""
    //hasExists := false
    //for i := len(ccs.ConfigSources) - 1; i >= 0; i-- {
    //    s := ccs.ConfigSources[i]
    //    v, err := s.Get(key)
    //    if err == nil {
    //        val = v
    //        hasExists = true
    //        break
    //    }
    //}
    //
    //if __reg.MatchString(val) {
    //    return ccs.evalValue(val)
    //}
    //if hasExists {
    //    return val, nil
    //} else {
    //    return val, errors.New("not exists for key: " + key)
    //}

    kv := ccs.KeyValue(key)

    return kv.value, kv.err
}

func (ccs *CompositeConfigSource) evalValue(value string) (string, error) {
    if strings.Contains(value, ccs.StartTag) {
        eval := fasttemplate.New(value, ccs.StartTag, ccs.EndTag)
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

    sub := __reg.FindStringSubmatch(value)
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
