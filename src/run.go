package main

import (
    "dockerEngine/src/container"
    "os"
    log "github.com/Sirupsen/logrus"
)

func Run(tty bool, command string) {
    parent := container.NewParentProcess(tty, command)
    if err := parent.Start(); err != nil {
        log.Error(err)
    }
    parent.Wait()
    os.Exit(-1)
}
