package ini

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"io/ioutil"
	"strings"
)

func ByIni(content string) *kvs.MapProperties {
	props, err := ReadIni(ioutil.NopCloser(strings.NewReader(content)))
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return &props.MapProperties
}
