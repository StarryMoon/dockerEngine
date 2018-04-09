package container

import (
    "syscall"
    "os/exec"
    "os"
    log "github.com/Sirupsen/logrus"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
    readPipe, writePipe, err := NewPipe()
    if err != nil {
        log.Errorf("New pipe error %v", err)
        return nil, nil
    }
//    args := []string{"init", command}
//    cmd := exec.Command("/proc/self/exe", args...)
    cmd := exec.Command("/proc/self/exe","init")
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWPID,
    }

    if tty {
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    }
    
    cmd.ExtraFiles = []*os.File{readPipe}
    cmd.Dir = "/root/busybox"
    return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
    read, write, err := os.Pipe()
    if err != nil {
        return nil, nil, nil 
    }
    return read, write, nil
}

