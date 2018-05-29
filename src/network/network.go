package netwoork

import (
    "github.com/vishvananda/netlink"
    "github.com/vishvananda/netns"
    "fmt"
    "net"
    "dockerEngine/src/container"
)

type Network struct {
    Name string                 //网络名
    IpRange *net.IPNet          //地址段
    Driver string               //网络驱动
}

type Endpoint struct {
    ID string `json:"id"`
    Device netlink.Veth `json:"dev"`
    IPAddress net.IP `json:"ip"`
    MacAddress net.HardwareAddr `json:"mac"`
    PortMapping []string `json:"portmapping"`
    Network  *Network
}

type NetworkDriver interface {
    Name() string
    Create(subnet string, name string) (*Network, error)
    Delete(network Network) error
    Connect(network *Network, endpoint *Endpoint) error
    Disconnect(network Network, endpoint *Endpoint) error
}

func CreateNetwork(driver string, subnet string, name string) error {
    _, cidr, _ := net.ParseCIDR(subnet)
    gatewayIP, err := ipAllocator.Allocate(cidr)
    if err != nil {
        return err
    }
    cidr.IP = gatewayIP
    fmt.Println("network/network.go CreateNetwork()  cidr.IP : ", cidr.IP)
    
    nw, err := drivers[driver].Create(cidr.String(), name)
    if err != nil {
        return err
    }

    return nw.dump(defaultNetworkPath)
}

func (nw *Network) dump(dumpPath string) error {
    if _, err := os.Stat(dumpPath); err != nil {
        if os.IsNoExist(err) {
            os.MkdirAll(dumpPath, 0644)
        }else {
            return err
        }
    }

    nwPath : = path.join(dumpPath, nw.Name)
    nwFile, err := os.OpenFile(nwPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
    if err != nil {
        logrus.Errorf("error: ", err)
        return err
    }
    defer nwFile.Close()
}
