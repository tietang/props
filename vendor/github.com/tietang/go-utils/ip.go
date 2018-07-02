package utils

import (
    "net"
    "strings"
)

//
func GetAllIP() []string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        panic(err)

    }
    size := len(addrs)
    ips := make([]string, size)
    for i, a := range addrs {

        //		println(a.Network(), a.String())
        ips[i] = strings.Split(a.String(), "/")[0]
        //		fmt.Println(strings.Split(a.String(), "/")[0])
        //		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
        //			if ipnet.IP.To4() != nil {
        //				//				os.Stdout.WriteString(ipnet.IP.String() + "\n")
        //				println(ipnet.IP.String())
        //			}
        //			println(ipnet.IP.String())
        //		}
    }
    return ips
}

var defaultExcludeFilters []string

func init() {
    defaultExcludeFilters = []string{"utun", "bridge", "docker", "vm", "veth", "br-", "vmnet"}
}
func GetExternalIP() (string, error) {
    ips, error := GetExternalIPs()
    return ips[0], error
}

func GetExternalIPs(excludeFilters ... string) ([]string, error) {
    filters := append(defaultExcludeFilters, excludeFilters...)

    ips := make([]string, 0)
    ifaces, err := net.Interfaces()
    if err != nil {
        return ips, err
    }

    for _, iface := range ifaces {
        if iface.Flags&net.FlagUp == 0 {
            continue // interface down
        }
        if iface.Flags&net.FlagLoopback != 0 {
            continue // loopback interface
        }
        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }
        isContinue := false
        for _, filter := range filters {
            if strings.Contains(iface.Name, filter) {
                isContinue = true
                break
            }
        }
        if isContinue {
            continue
        }

        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            if ip == nil || ip.IsLoopback() {
                continue
            }
            ip = ip.To4()
            if ip == nil {
                continue // not an ipv4 address
            }
            ips = append(ips, ip.String())
        }
    }
    return ips, nil
    //, errors.New("are you connected to the network?")
}
