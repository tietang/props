package props

import (
    "net/http"
    "time"
)

var testEtcd *MockTestEtcd

type MockTestEtcd struct {
    Address string
}

func GetOrNewMockTestEtcd(address string) *MockTestEtcd {
    if testEtcd == nil {
        testEtcd = &MockTestEtcd{
            Address: address,
        }
    }

    return testEtcd
}

func (m *MockTestEtcd) StartMockEtcd() <-chan int {
    ec := make(chan int, 1)
    isStarted := m.CheckEtcdIsStarted()

    if isStarted {
        ec <- 1
        return ec
    }
    command := "etcd"
    params := []string{"agent", "-dev"}
    started := execCommand(command, params)

    if started {
        ec <- 1
    } else {
        ec <- 0
    }
    return ec
}

func (m *MockTestEtcd) WaitingForEtcdStarted() {
    for {
        isStarted := m.CheckEtcdIsStarted()
        if isStarted {
            break
        }
        time.Sleep(time.Millisecond * 200)
    }
}

func (m *MockTestEtcd) CheckEtcdIsStarted() bool {
    res, err := http.Get("http://" + m.Address)
    if err != nil {
        return false
    }
    if res != nil && res.StatusCode == 200 {
        res.Body.Close()
        return true
    }
    return false
}
