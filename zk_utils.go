package props

import (
    "github.com/samuel/go-zookeeper/zk"
    "path"
    "log"
    "fmt"
    "bufio"
    "io"
    "os/exec"
    "os"
)

var flags = int32(0)
var acl = zk.WorldACL(zk.PermAll)

func ZkCreateString(conn *zk.Conn, path string, value string) (string, error) {
    return ZkCreate(conn, path, []byte(value))
}

func ZkCreate(conn *zk.Conn, nodePath string, value []byte) (string, error) {

    d, _ := path.Split(nodePath)
    ppath := path.Clean(d)
    if !ZkExits(conn, ppath) {
        ZkCreate(conn, ppath, []byte(""))
    }
    return conn.Create(nodePath, []byte(value), flags, acl)
}

func ZkExits(conn *zk.Conn, path string) bool {
    b, _, e := conn.Exists(path)
    if e != nil {
        return false
    }
    return b
}

func ZkWatchNodeCreated(conn *zk.Conn, path string) {
    log.Println("watchNodeCreated")
    for {
        _, _, ch, _ := conn.ExistsW(path)
        e := <-ch
        log.Println("ExistsW:", e.Type, "Event:", e)
        if e.Type == zk.EventNodeCreated {
            log.Println("NodeCreated ")
            return
        }
    }
}
func ZkWatchNodeDeleted(conn *zk.Conn, path string) {
    log.Println("watchNodeDeleted")
    for {
        _, _, ch, _ := conn.ExistsW(path)
        e := <-ch
        log.Println("ExistsW:", e.Type, "Event:", e)
        if e.Type == zk.EventNodeDeleted {
            log.Println("NodeDeleted ")
            return
        }
    }
}

func ZkWatchNodeDataChange(conn *zk.Conn, path string) {
    for {
        _, _, ch, _ := conn.GetW(path)
        e := <-ch
        log.Println("GetW('"+path+"'):", e.Type, "Event:", e)
    }
}

func ZkWatchChildrenChanged(conn *zk.Conn, path string) {
    for {
        _, _, ch, _ := conn.ChildrenW(path)
        e := <-ch
        log.Println("ChildrenW:", e.Type, "Event:", e)
    }
}

func execCommand(commandName string, params []string) bool {

    cmd := exec.Command(commandName, params...)

    //显示运行的命令
    fmt.Println(cmd.Args)

    stdout, err := cmd.StdoutPipe()

    if err != nil {
        fmt.Println(err)
        return false
    }

    cmd.Start()

    reader := bufio.NewReader(stdout)

    //实时循环读取输出流中的一行内容
    for {
        line, err2 := reader.ReadString('\n')
        if err2 != nil || io.EOF == err2 {
            fmt.Println("exit consul start.")
            break
        }
        fmt.Print(line)
    }
    cmd.Process.Signal(os.Kill)
    cmd.Wait()
    return true
}
