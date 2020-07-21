package yam

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"io"
	"os"
	"path"
)

type YamlConfigSource struct {
	YamlProperties
	name     string
	fileName string
}

func NewYamlConfigSource(fileName string) *YamlConfigSource {
	name := path.Base(fileName)
	return NewYamlConfigSourceByFile(name, fileName)
}

func NewYamlConfigSourceByFile(name, file string) *YamlConfigSource {

	f, err := os.Open(file)
	if err != nil {
		log.Error(err)
		return nil
	}
	defer f.Close()
	s := NewYamlConfigSourceByReader(name, f)
	s.fileName = file
	return s
}

func NewYamlConfigSourceByReader(name string, r io.Reader) *YamlConfigSource {

	s := &YamlConfigSource{}
	s.name = name
	s.Values = make(map[string]string)
	s.fileName = "no-file"
	if s.Values == nil {
		s.Values = make(map[string]string)
	}
	s.Load(r)
	return s
}

func NewIniFileCompositeConfigSource(fileNames ...string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "yamlFiles"
	for _, file := range fileNames {
		c := NewYamlConfigSource(file)
		s.Add(c)
	}
	return s
}

func (s *YamlConfigSource) Name() string {
	return s.name
}

func (s *YamlConfigSource) FileName() string {
	return s.fileName
}
