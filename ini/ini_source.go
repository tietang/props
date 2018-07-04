package ini

import (
    log "github.com/sirupsen/logrus"
    "path"
    "path/filepath"
    "io"
    "github.com/tietang/props/kvs"
)

const (
    KEY_INI_CURRENT_DIR = "ini.current.dir"
)

//ini 文件支持
type IniFileConfigSource struct {
    kvs.MapProperties
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
        m = p.Values
    } else {
        log.WithField("error", err.Error()).Fatal("read file: ")
    }
    s := &IniFileConfigSource{}
    s.name = name
    s.Values = m
    s.fileName = file
    if s.Values == nil {
        s.Values = make(map[string]string)
    }
    s.Set(KEY_INI_CURRENT_DIR, filepath.Dir(file))
    return s
}

func NewIniFileConfigSourceByReader(name string, r io.Reader) *IniFileConfigSource {
    p, err := ReadIni(r)
    var m map[string]string
    if err == nil {
        m = p.Values
    } else {
        log.WithField("error", err.Error()).Fatal("read file: ")
    }
    s := &IniFileConfigSource{}
    s.name = name
    s.Values = m
    s.fileName = "no-file"
    return s
}

func NewIniFileCompositeConfigSource(fileNames ...string) *kvs.CompositeConfigSource {
    s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
    s.ConfName = "iniFiles"
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
