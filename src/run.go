package main

import (
    "dockerEngine/src/container"
    "dockerEngine/src/cgroups"
    "dockerEngine/src/cgroups/subsystems"
    "os"
    log "github.com/Sirupsen/logrus"
)

func RunCmd(tty bool, command string, res *subsystems.ResourceConfig) {
    parent := container.NewParentProcess(tty, command)
    if err := parent.Start(); err != nil {
        log.Error(err)
    }

    // create cgroup
    cgroupManager := cgroups.NewCgroupManager("dockerEngine-cgroup")
    defer cgroupManager.Destory()
    cgroupManager.Set(res)
    cgroupManager.Apply(parent.Process.Pid)

    parent.Wait()
    os.Exit(-1)
}
