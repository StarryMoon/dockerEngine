package netwoork

import (
    "github.com/vishvananda/netlink"
    "github.com/vishvananda/netns"
    "fmt"
    "net"
    "dockerEngine/src/container"
)

var (
    defaultNetworkPath = "/var/run/dockerEngine/network/network/"
)

type NetworkSeg struct {            //建立子网
    Name string                 //网络名
    IpRange *net.IPNet          //地址段
    Driver string               //网络驱动
}

type Endpoint struct {         //veth
    ID string `json:"id"`
    Device netlink.Veth `json:"dev"`
    IPAddress net.IP `json:"ip"`
    MacAddress net.HardwareAddr `json:"mac"`
    PortMapping []string `json:"portmapping"`
    Network  *NetworkSeg
}

type NetworkDriver interface {
    Name() string
    Create(subnet string, name string) (*NetworkSeg, error)
    Delete(network NetworkSeg) error
    Connect(network *NetworkSeg, endpoint *Endpoint) error
    Disconnect(network NetworkSeg, endpoint *Endpoint) error
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

//保存网络信息
func (nw *NetworkSeg) dump(dumpPath string) error {
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

    nwJson, err := json.Marshal(nw)
    if err != nil {
        logrus.Errorf("error: ", err)
        return err
    }

    _, errr := nwFile.Write(nwJson)
    if err != nil {
        logrus.Errorf("error: ", err)
        return err
    }

    return nil
}

func (nw *NetworkSeg) load(dumpPath string) error {
//    nwPath := path.join(dumpPath, nw.Name)
    nwConfigFile, err := os.Open(dumpPath)
    defer nwConfigFile.Close()
    if err != nil {
        return err
    }

    nwJson := make([]byte, 2000)
    n, err := nwConfigFile.Read(nwJson)
    if err != nil {
        return err
    }

    err = json.Unmarshal(nwJson[:n], nw)
    if err != nil {
        logrus.Errorf("Error load nw info", err)
        return err
    }

    return nil
}

func Connect(networkName string, cinfo *container.ContainerInfo) error {
    network, ok := networks[networkName]
    if !ok {
        return fmt.Errorf("No such network: %s", networkName)
    }

    ip, err := ipAllocator.Allocate(network.IPRange)
    if err != nil {
        return err
    }

    ep := &Endpoint{
        ID: fmt.Sprintf("%s-%s", cinfo.Id, networkName),
        IPAddress: ip,
        Network: network,
        PortMapping: cinfo.PortMapping,
    }

    if err = drivers[network.Driver].Connect(network, ep); err != nil {
        return err
    }

    if err = configEndpointIpAddressAndRoute(ep, cinfo); err != nil {
        return err
    }

    return configPortMapping(ep, cinfo)
}
