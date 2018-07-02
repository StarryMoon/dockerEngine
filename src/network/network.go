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
    networks = map[string]*NetworkSeg{}
    drivers = map[string]NetworkDriver{}
)

type NetworkSeg struct {        //建立子网
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

func Init() error {
    //加载网络驱动
    var bridgeDriver = BridgeNetworkDriver()
    drivers[bridgeDriver.Name()] = &bridgeDriver
    if _, err := os.Stat(defaultNetworkPath); err != nil {
        if os.IsNoExit(err) {
            os.MkdirAll(defaultNetworkPath, 0644)
        } else {
            return err
        }
    }
    
    filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
        //自动拼凑 nwPath = defaultNetworkPath + fileName
        if info.IsDir() {
            return nil
        }

        _, nwName := path.Split(nwPath)
        nw := &Network{
            Name: nwName,
        }

        if err := nw.load(nwPath); err != nil {
            logrus.Errorf("error load network: %s", err)
        }

        networks[nwName] = nw
        return nil
    })
    return nil
}

func ListNetwork() {
    w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
    fmt.Fprint(w, "NAME\tIPRange\tDriver\n")
    for _, nw := range networks {
        fmt.Fprintf(w, "%s\t%s\t%s\n",
            nw.Name,
            nw.IPRange.String(),
            nw.Driver,
        )
    }

    if err := w.Flush(); err != nil {
        logrus.Errorf("Flush error %v", err)
        return
    }
}

func DeleteNetwork(networkName string) error {
    nw, ok := networks[networkName]
    if !ok {
        return fmt.Errorf("No such network: %s", networkName)
    }

    if err := ipAllocator.Release(nw.IPRange, &nw.IPRange.IP); err != nil {
        retrun fmt.Errorf("Error Remove Network gateway ip: %s", err)
    }

    if err := drivers[nw.Driver].Delete(*nw); err != nil {
        retrun fmt.Errorf("Error Remove Network DriveError: %s", err)
    }

    return nw.remove(defaultNetworkPath)
}

//删除文件
func (nw *Network) remove(dumpPath string) error {
    if _, err := os.Stat(path.Join(dumpPath, nw.Name)); err != nil {
        if os.IsNoExist(err) {
             return nil
        } else {
            return err 
        }
    } else {
        return  os.Remove(path.Join(dumpPath, nw.Name))
    }
}
