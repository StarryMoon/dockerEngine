package network 

import (
    "net"
    "os"
)

const ipamDefaultAllocatorPath = "/var/run/dockerEngine/network/ipam/subnet.json"

type IPAM struct {
    SubnetAllocatorPath string
    Subnets *map[string]string
}

//全局变量
var ipAllocator = &IPAM{
    SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

func (ipam *IPAM) load() error {
    if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
        if os.IsNotExist(err) {
            return nil
        } else {
            return err
        }
    }

    subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
    defer subnetConfigFile.Close()
    if err != nil {
        return err
    }

    subnetJson := make([]byte, 2000)
    n, err := subnetConfigFile.Read(subnetJson)
    if err != nil {
        return err
    }
    
    err = json.Unmarshal(subnetJson[:n], ipam.Subnets)
    if err != nil {
        log.Errorf("Error dump allocation info, %v", err)
        return err
    }

    return nil
}

func (ipam *IPAM) dump() error {
    ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
    if _, err := os.Stat(ipamConfigFileDir); err != nil {
        if os.IsNotExist(err) {
             os.MkdirAll(ipamConfigFileDir, 0644)
        }else {
            return err
        }
    }

    subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC | os.O_WRONLY | OS.O_CREATE, 0644)
    defer subnetConfigFile.close()
    if err != nil {
        return err
    }

    ipamConfigJson, err := json.Marshal(ipam.Subnets)
    if err != nil {
        return err
    }

    _, err = subnetConfigFile.Write(ipamConfigJson)
    if err != nil {
        return err
    }

    return nil
}
