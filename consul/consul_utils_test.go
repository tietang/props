package consul

import (
	"github.com/tietang/props/v3/kvs"
	"net/http"
	"time"
)

var consul_mock_started = false

func init() {
	//address := "172.16.1.248:8500"
	address := "127.0.0.1:8500"
	GetOrNewMockTestConsul(address)
	go kvs.ExecCommand("pwd", "-LP")
	if !consul_mock_started {
		go testConsul.StartMockConsul()
	}
	testConsul.WaitingForConsulStarted()
}

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
	started := kvs.ExecCommand(command, params...)

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
