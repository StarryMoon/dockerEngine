package main


import (
    "fmt"
    log "github.com/Sirupsen/logrus"
    "github.com/StarryMoon/dockerEngine/container"
    "io/ioutil"
    "encoding/json"
    "strings"
    "os/exec"
    "os"
    _ "github.com/StarryMoon/dockerEngine/nsenter"
)

func GetContainerPidByName(containerName string) (string, error) {
    dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
    configFilePath := dirURL + container.ConfigName
    contentBytes, err := ioutil.ReadFile(configFilePath)
    if err != nil {
        return "", err
    }
    var containerInfo container.ContainerInfo
    if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
        return "", err
    }
    return containerInfo.Pid, nil
}
