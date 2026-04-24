package tom

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
)

func ByToml(content string) *kvs.MapProperties {
	y := NewTomlProperties()
	err := y.Load(strings.NewReader(content))
	if err != nil {
		log.Error(err)
		return nil
	}
	return &y.MapProperties
}
