package zk

import "github.com/tietang/props/kvs"

var zk_mock_started bool = false

func init() {
    if !zk_mock_started {
        go kvs.ExecCommand("pwd", "-LP")
        go StartMockTestZookeeper()
    }
}
func StartMockTestZookeeper() <-chan int {
    ec := make(chan int, 1)
    if !zk_mock_started {

        command := "java"
        params := []string{"-jar", "zookeeper/mock.jar"}
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
