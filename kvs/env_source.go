package kvs

import (
    "os"
    "strings"
)

//内部使用不扩展
type envConfigSource struct {
    MapProperties
    name string
}

func newEnvConfigSource() *envConfigSource {
    p := &envConfigSource{}
    p.Values = make(map[string]string)
    p.Init()
    return p
}

func (e *envConfigSource) Init() {
    envs := os.Environ()
    for _, v := range envs {
        idx := strings.Index(v, "=")
        if idx > 0 {
            e.Set(v[:idx], v[idx+1:])
        }
    }
}
func (e *envConfigSource) Name() string {

    return e.name
}
