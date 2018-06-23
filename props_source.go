package props

import (
    log "github.com/sirupsen/logrus"
    "path"
    "path/filepath"
)

const (
    KEY_PROPS_CURRENT_DIR = "current.dir"
)

type PropertiesConfigSource struct {
    MapProperties
    name     string
    fileName string
}

func NewPropertiesConfigSource(fileName string) *PropertiesConfigSource {
    name := path.Base(fileName)
    return NewPropertiesConfigSourceByFile(name, fileName)
}

func NewPropertiesConfigSourceByFile(name, file string) *PropertiesConfigSource {

    p, err := ReadPropertyFile(file)

    var m map[string]string
    if err == nil {
        m = p.values
    } else {
        log.WithField("error", err.Error()).Error("read file: ")
    }
    s := &PropertiesConfigSource{}
    s.name = name
    s.values = m
    s.fileName = file
    s.Set(KEY_PROPS_CURRENT_DIR, filepath.Dir(file))
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

func NewPropertiesCompositeConfigSource(fileNames ...string) *CompositeConfigSource {
    s := NewEmptyNoSystemEnvCompositeConfigSource()
    s.name = "properties"
    for _, file := range fileNames {
        c := NewPropertiesConfigSource(file)
        s.Add(c)
    }
    return s
}
func NewEmptyMapConfigSource(name string) *PropertiesConfigSource {
    s := &PropertiesConfigSource{}
    if name == "" {
        s.name = "Map"
    } else {
        s.name = name
    }

    s.values = make(map[string]string)
    return s
}
func (s *PropertiesConfigSource) Name() string {
    return s.name
}

func (s *PropertiesConfigSource) FileName() string {
    return s.fileName
}
