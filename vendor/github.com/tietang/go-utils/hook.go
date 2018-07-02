package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	// or "runtime"
)

type signalFunc func(s os.Signal, arg interface{})

type Hook struct {
	m map[os.Signal]signalFunc
}

func NewHook() *Hook {
	ss := new(Hook)
	ss.m = make(map[os.Signal]signalFunc)
	return ss
}

func (set *Hook) Register(s os.Signal, handler signalFunc) {
	if _, found := set.m[s]; !found {
		set.m[s] = handler
	}
}

func (set *Hook) Handle(sig os.Signal, arg interface{}) (err error) {
	if _, found := set.m[sig]; found {
		set.m[sig](sig, arg)
		return nil
	} else {
		return fmt.Errorf("No handler available for signal %v", sig)
	}

	panic("won't reach here")
}

func demo() {
	hook := NewHook()
	handler := func(s os.Signal, arg interface{}) {
		fmt.Printf("handle signal: %v\n", s)
	}
	//SIGINT,SIGTERM,SIGQUIT
	//Interrupt Signal = syscall.SIGINT interrupt
	//Kill      Signal = syscall.SIGKILL killed
	//syscall.SIGTERM terminated

	hook.Register(os.Interrupt, handler)
	hook.Register(os.Kill, handler)
	hook.Register(syscall.SIGTERM, handler)

	for {
		c := make(chan os.Signal)
		//		var sigs []os.Signal
		//		for sig := range ss.m {
		//			sigs = append(sigs, sig)
		//		}
		signal.Notify(c)
		sig := <-c

		err := hook.Handle(sig, nil)
		if err != nil {
			fmt.Printf("unknown signal received: %v\n", sig)
			//			os.Exit(1)
		}
	}
}
