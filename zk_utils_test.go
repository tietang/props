package props

var zk_mock_started bool = false

func StartMockTestZookeeper() <-chan int {
    ec := make(chan int, 1)
    if !zk_mock_started {

        command := "java"
        params := []string{"-jar", "zookeeper/mock.jar"}
        started := execCommand(command, params)

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
