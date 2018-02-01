package container

import (
    "os"
    "syscall"
    log "github.com/Sirupsen/logrus"
    "fmt"
    "io/ioutil"
    "strings"
    "os/exec"
)

func RunContainerInitProcess() error {
    cmdArray := readUserCommand()
    if cmdArray == nil || len(cmdArray) == 0 {
        return fmt.Errorf("Run container get user command error, cmdArray is nil")
    }
    
    log.Infof("command %s", cmdArray)
    fmt.Println("command : ", cmdArray)

    //setUpMount()

    path, err := exec.LookPath(cmdArray[0])
    if err != nil {
        log.Errorf("Exec loop path error %v", err)
        return err
    }
    log.Infof("Find path %s", path)
    
    defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
    syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
//    argv := []string{command}
    if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
        log.Errorf(err.Error())
    }
   
    return nil
}

func readUserCommand() []string {
    pipe := os.NewFile(uintptr(3), "pipe")
    msg, err := ioutil.ReadAll(pipe)
    if err != nil {
        log.Errorf("init command read pipe error %v", err)
        return nil
    }
    msgStr := string(msg)
    return strings.Split(msgStr, " ")
}
