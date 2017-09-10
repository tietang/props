package props

import (
    log "github.com/sirupsen/logrus"
    "path"
    "path/filepath"
)

const (
    KEY_INI_CURRENT_DIR = "ini.current.dir"
)

//ini 文件支持
type IniFileConfigSource struct {
    MapProperties
    name     string
    fileName string
}

func NewIniFileConfigSource(fileName string) *IniFileConfigSource {
    name := path.Base(fileName)
    return NewIniFileConfigSourceByFile(name, fileName)
}

func NewIniFileConfigSourceByFile(name, file string) *IniFileConfigSource {

    p, err := ReadIniFile(file)

    var m map[string]string
    if err == nil {
        m = p.values
    } else {
        log.WithField("error", err.Error()).Info("read file: ")
    }
    s := &IniFileConfigSource{}
    s.name = name
    s.values = m
    s.fileName = file
    s.Set(KEY_INI_CURRENT_DIR, filepath.Dir(file))
    return s
}

func NewIniFileCompositeConfigSource(fileNames ...string) *CompositeConfigSource {
    s := NewEmptyCompositeConfigSource()
    s.name = "iniFiles"
    for _, file := range fileNames {
        c := NewIniFileConfigSource(file)
        s.Add(c)
    }
    return s
}

func (s *IniFileConfigSource) Name() string {
    return s.name
}

func (s *IniFileConfigSource) FileName() string {
    return s.fileName
}
