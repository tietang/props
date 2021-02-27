package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
)

func main() {
	//获取程序运行文件所在的路径
	file := kvs.GetCurrentFilePath("api.ini", 2)
	log.Info("config file: ", file)
}
