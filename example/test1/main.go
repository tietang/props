package main

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp/v3"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"gopkg.in/yaml.v3"
)

func main() {
	var aiConfig = &AIConfig{}
	aiConfigFile := kvs.GetCurrentFilePath("ai.yaml", 2)
	//log.Info("aiConfigFile:%v", aiConfigFile)
	//aiConfigSource := yam.NewYamlConfigSource(aiConfigFile)
	//err := aiConfigSource.Unmarshal(aiConfig)
	//if err != nil {
	//	log.Error("解析配置文件失败:", err)
	//}

	data, err := os.ReadFile(aiConfigFile)
	if err != nil {
		log.Error("读取配置文件失败:", err)
	}
	err = yaml.Unmarshal(data, aiConfig)
	if err != nil {
		log.Error("解析配置文件失败:", err)
	}
	pp.Print(aiConfig)
	fmt.Printf("%+v\n", aiConfig)
}
