package props

import (
	log "github.com/sirupsen/logrus"
	"path"
)

type PropertiesConfigSource struct {
	MapProperties
	name string
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
		log.WithField("error", err.Error()).Info("read file: ")
	}
	s := &PropertiesConfigSource{}
	s.name = name
	s.values = m
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

func (s *PropertiesConfigSource) Name() string {
	return s.name
}
