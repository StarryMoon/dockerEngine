package container

import (
    "os"
    "syscall"
    log "github.com/Sirupsen/logrus"
    "fmt"
)

func RunContainerInitProcess(command string, args []string) error {
    log.Infof("command %s", command)
    fmt.Println("command : ", command)   

    defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NSUID | syscall.MS_NODEV
    syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
    argv := []string{command}
    if err := syscall.Exec(command, argv, os.Environ()); err != nil {
        log.Errorf(err.Error())
    }
    
    retrun nil
}
