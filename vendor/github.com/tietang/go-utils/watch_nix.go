// +build darwin linux

package utils

import (
    "os"
    "os/signal"
    //"syscall"
    "syscall"
)

//
func Notify(callback func()) {
    //热更新配置可能有多种触发方式，这里使用系统信号量sigusr1实现
    //通过kill -USR1 PID|kill -30 PID  来更新
    sigs := make(chan os.Signal, 1)

    signal.Notify(sigs, syscall.SIGUSR1) //30
    //signal.Notify(sigs, syscall.Signal(0xa))
    //go func() {
    //    for syscall.Signal(0xa) == <-sigs {
    //        log.Print("Recieved 0xa, reloading config")
    //        callback()
    //    }
    //}()
    go func() {
        for {
            <-sigs
            callback()
        }
    }()

}
