package container

import (
    "os"
    "syscall"
    "os/exec"
)

func NewParentProcess(tty bool, command string) *exec.Command {
    args := []string{"init", command}
    cmd := exec.Command("/proc/self/exe", args...)
    cmd.SysProcAttr = &syscall.SysProcAttr{
        CloneFlags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS | syscall.CLONE_NET,
    }

    if tty {
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    }
    return cmd
}
