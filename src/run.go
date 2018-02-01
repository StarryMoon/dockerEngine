package main

import (
    "dockerEngine/src/container"
    "dockerEngine/src/cgroups"
    "dockerEngine/src/cgroups/subsystems"
    "os"
    "strings"
    log "github.com/Sirupsen/logrus"
)

func RunCmd(tty bool, command []string, res *subsystems.ResourceConfig) {
    parent, writePipe := container.NewParentProcess(tty)
    if parent == nil {
        log.Errorf("New parent process error")
        return
    }
    if err := parent.Start(); err != nil {
        log.Error(err)
    }

    // create cgroup
    cgroupManager := cgroups.NewCgroupManager("dockerEngine-cgroup")
    defer cgroupManager.Destory()
    cgroupManager.Set(res)
    cgroupManager.Apply(parent.Process.Pid)

    sendInitCommand(command, writePipe)
    parent.Wait()
    os.Exit(-1)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
    command := strings.Join(cmdArray, " ")
    log.Infof("command all is %s ", command)
    writePipe.WriteString(command)
    writePipe.Close()
}
