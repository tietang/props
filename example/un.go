package main

import (
    "fmt"
    "github.com/tietang/props"
    "time"
)

type Port struct {
    Port    int  `val:"8080"`
    Enabled bool `val:"true"`
}
type ServerProperties struct {
    _prefix string        `prefix:"http.server"`
    Port    Port
    Timeout int           `val:"1"`
    Enabled bool
    Foo     int           `val:"1"`
    Time    time.Duration `val:"1s"`
    Float   float32       `val:"0.000001"`
    Params  map[string]string
    Times      map[string]time.Duration
}

func main() {
    //t := ServerProperties{}
    //s := reflect.ValueOf(&t).Elem()
    //typeOfT := s.Type()
    //
    //for i := 0; i < s.NumField(); i++ {
    //	f := s.Field(i)
    //	fmt.Printf("%d: %s %s = %v\n", i,
    //		typeOfT.Field(i).Name, f.Type(), f.Interface())
    //}
    //
    //s.Field(0).SetInt(25)
    //s.Field(1).SetString("nicky")
    //fmt.Println(t)

    //
    p := props.NewMapProperties()
    p.Set("http.server.port.port", "8080")
    p.Set("http.server.params.k1", "v1")
    p.Set("http.server.params.k2", "v2")
    p.Set("http.server.Times.m1", "1s")
    p.Set("http.server.Times.m2", "1h")
    p.Set("http.server.Times.m3", "1us")
    p.Set("http.server.port.enabled", "false")
    p.Set("http.server.timeout", "1234")
    p.Set("http.server.enabled", "true")
    p.Set("http.server.time", "10s")
    p.Set("http.server.float", "23.45")
    p.Set("http.server.foo", "23")
    s := &ServerProperties{
        Foo:   1234,
        Float: 1234.5,
    }
    p.Unmarshal(s)
    fmt.Println(s)

}
