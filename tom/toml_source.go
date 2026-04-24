package tom

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
)

var _ kvs.ConfigSource = new(TomlConfigSource)

type TomlConfigSource struct {
	TomlProperties
	name     string
	fileName string
}

func NewTomlConfigSource(fileName string) *TomlConfigSource {
	name := filepath.Base(fileName)
	return NewTomlConfigSourceByFile(name, fileName)
}

func NewTomlConfigSourceByFile(name, file string) *TomlConfigSource {

	f, err := os.Open(file)
	if err != nil {
		log.Error(err)
		return nil
	}
	defer f.Close()
	s := NewTomlConfigSourceByReader(name, f)
	s.fileName = file
	return s
}

func NewTomlConfigSourceByReader(name string, r io.Reader) *TomlConfigSource {

	s := &TomlConfigSource{}
	s.name = name
	s.Values = make(map[string]string)
	s.fileName = "no.txt-file"
	if s.Values == nil {
		s.Values = make(map[string]string)
	}
	s.Load(r)
	return s
}

func NewTomlFileCompositeConfigSource(fileNames ...string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "TomlFiles"
	for _, file := range fileNames {
		c := NewTomlConfigSource(file)
		s.Add(c)
	}
	return s
}

func (s *TomlConfigSource) Name() string {
	return s.name
}

func (s *TomlConfigSource) FileName() string {
	return s.fileName
}
