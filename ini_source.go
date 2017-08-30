package props

import (
    log "github.com/sirupsen/logrus"
    "path"
    "path/filepath"
)

const (
    KEY_INI_CURRENT_DIR = "ini.current.dir"
)

type IniConfigSource struct {
    MapProperties
    name     string
    fileName string
}

func NewIniConfigSource(fileName string) *IniConfigSource {
    name := path.Base(fileName)
    return NewIniConfigSourceByFile(name, fileName)
}

func NewIniConfigSourceByFile(name, file string) *IniConfigSource {

    p, err := ReadIniFile(file)

    var m map[string]string
    if err == nil {
        m = p.values
    } else {
        log.WithField("error", err.Error()).Info("read file: ")
    }
    s := &IniConfigSource{}
    s.name = name
    s.values = m
    s.fileName = file
    s.Set(KEY_INI_CURRENT_DIR, filepath.Dir(file))
    return s
}

func NewIniConfigSourceByMap(name string, kv map[string]string) *IniConfigSource {
    s := &IniConfigSource{}
    s.name = name
    if kv == nil {
        s.values = make(map[string]string)
    } else {
        s.values = kv
    }
    return s
}

func (s *IniConfigSource) Name() string {
    return s.name
}

func (s *IniConfigSource) FileName() string {
    return s.fileName
}
