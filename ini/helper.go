package ini

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
)

func ByIni(content string) *kvs.MapProperties {
	props, err := ReadIni(io.NopCloser(strings.NewReader(content)))
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return &props.MapProperties
}
