package netwoork

import (
    "github.com/vishvananda/netlink"
    "github.com/vishvananda/netns"
    "fmt"
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
