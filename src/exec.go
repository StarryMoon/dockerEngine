package main

import (
    log "github.com/Sirupsen/logrus"
    "fmt"
    "os"
    "io/ioutil"
    "encoding/json"
    "strings"
    "os/exec"
    "dockerEngine/src/container"
    _ "dockerEngine/src/nsenter"
)

const ENV_EXEC_PID = "dockerEngine_pid"
const ENV_EXEC_CMD = "dockerEngine_cmd"

func execContainer(containerName string, comArray []string) {
    pid, err := getContainerPidByName(containerName)
    if err != nil {
        log.Errorf("Exec container getContainerPidByName %s error %v", containerName, err)
        return
    }

    cmdStr := strings.Join(comArray, " ")
    log.Infof("container pid %s", pid)
    log.Infof("command %s", cmdStr)

    cmd := exec.Command("/proc/self/exe", "exec")     //trigger the execution of exec again
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    os.Setenv(ENV_EXEC_PID, pid)
    os.Setenv(ENV_EXEC_CMD, cmdStr)

    if err := cmd.Run(); err != nil {
        log.Errorf("Exec container %s error %v", containerName, err)
    }
}


func getContainerPidByName(containerName string) (string, error) {
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

    return containerInfo.Pid,nil
}
