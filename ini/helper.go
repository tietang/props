package ini

import (
	"github.com/prometheus/common/log"
	"github.com/tietang/props/kvs"
	"strings"
)

func ByIni(content string) *kvs.MapProperties {
	props, err := ReadIni(strings.NewReader(content))
	if err != nil {
		log.Error(err)
		return nil
	}
	return &props.MapProperties
}
