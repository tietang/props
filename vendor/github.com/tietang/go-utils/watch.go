package utils

import (
    "os"
    "os/signal"
    //"syscall"
)

func Watch(sig os.Signal, callback func()) {
    //热更新配置可能有多种触发方式，这里使用系统信号量sigusr1实现
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, sig)
    go func() {
        for {
            <-sigs
            callback()
        }
    }()

}
