package props

import (
    "net/http"
    "time"
)

var testConsul *MockTestConsul

type MockTestConsul struct {
    Address string
}

func GetOrNewMockTestConsul(address string) *MockTestConsul {
    if testConsul == nil {
        testConsul = &MockTestConsul{
            Address: address,
        }
    }

    return testConsul
}
func (m *MockTestConsul) StartMockConsul() <-chan int {
    ec := make(chan int, 1)
    isStarted := m.CheckConsulIsStarted()

    if isStarted {
        ec <- 1
        return ec
    }
    command := "consul"
    params := []string{"agent", "-dev"}
    started := execCommand(command, params)

    if started {
        ec <- 1
    } else {
        ec <- 0
    }
    return ec
}

func (m *MockTestConsul) WaitingForConsulStarted() {
    for {
        isStarted := m.CheckConsulIsStarted()
        if isStarted {
            break
        }
        time.Sleep(time.Millisecond * 200)
    }
}

func (m *MockTestConsul) CheckConsulIsStarted() bool {
    res, err := http.Get("http://"+m.Address)
    if err != nil {
        return false
    }
    if res != nil && res.StatusCode == 200 {
        res.Body.Close()
        return true
    }
    return false
}
