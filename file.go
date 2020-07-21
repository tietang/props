package props

import (
	"errors"
	"github.com/tietang/props/v3/ini"
	"github.com/tietang/props/v3/kvs"
	"path/filepath"
	"strings"
)

const (
	INI_FILE_EXT   = ".ini,.conf,.cfg,.flowconfig,.config"
	PROPS_FILE_EXT = ".props,.properties,.messages"
)

func ReadFile(f string) (kvs.ConfigSource, error) {
	ext := filepath.Ext(f)
	if strings.Contains(INI_FILE_EXT, ext) {
		return ini.ReadIniFile(f)
	}
	if strings.Contains(PROPS_FILE_EXT, ext) {
		return kvs.ReadPropertyFile(f)
	}
	return nil, errors.New("Unsupported file type: " + ext)
}
