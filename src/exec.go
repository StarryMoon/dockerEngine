package main

import (
    log "github.com/Sirupsen/logrus"
    "fmt"
    "os"
    "io/ioutil"
    "encoding/json"
    "text/tabwriter"
//    "os/exec"
    "dockerEngine/src/container"
)

func ExecContainer(icontainerName string, comArray []string) {
    pid, err := getContainerPidByName(containerName)
    if err != nil {
        log.Errorf("Exec container getContainerPidByName %s error %v", containerName, err)
        return
    }

    cmdStr := string.Join(comArray, " ")
    log.Infof("container pid %s", pid)
    log.Infof("command %s", cmdStr)


}
