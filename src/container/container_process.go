package container

import (
    "syscall"
    "os/exec"
    "os"
)

func NewParentProcess(tty bool, command string) {
    args := []string("init", command)
    cmd := exec.command("/proc/self/exe", args...)
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID
    }

    if tty {
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    }

    return cmd
}