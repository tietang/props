package zk

import (
	"fmt"
	"github.com/tietang/props/v3/kvs"
	"os"
	"path/filepath"
)

var zk_mock_started bool = false

func init() {
	if !zk_mock_started {
		go kvs.ExecCommand("pwd", "-LP")
		go StartMockTestZookeeper()
	}
	fmt.Println(os.Getwd())
}

func StartMockTestZookeeper() <-chan int {
	ec := make(chan int, 1)
	pwd, _ := os.Getwd()
	jar := filepath.Join(pwd, "zookeeper/mock.jar")
	fmt.Println(jar)
	if !zk_mock_started {

		command := "java"
		params := []string{"-jar", jar}
		started := kvs.ExecCommand(command, params...)

		if started {
			ec <- 1
		} else {
			ec <- 0
		}
	} else {
		ec <- 1
	}

	return ec
}
