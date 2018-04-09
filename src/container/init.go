package container

import (
    "os"
    "syscall"
    log "github.com/Sirupsen/logrus"
    "fmt"
    "path/filepath"
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

    setUpMount()

    path, err := exec.LookPath(cmdArray[0])
    if err != nil {
        log.Errorf("Exec loop path error %v", err)
        return err
    }
    log.Infof("Find path %s", path)
    
//    defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
//    syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
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

func pivotRoot(root string) error {

    //重新挂载root文件系统，确保不在同一颗目录树
    if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
        return fmt.Errorf("Mount rootfs to itself error: %v", err)
    }

    pivotDir := filepath.Join(root, ".pivot_root")
    if err := os.Mkdir(pivotDir, 0777); err != nil {
        return err
    }

    if err := syscall.PivotRoot(root, pivotDir); err != nil {
        return fmt.Errorf("pivot_root %v", err) 
    }

    if err := syscall.Chdir("/"); err != nil {
        return fmt.Errorf("chdir / %v", err)
    }

    pivotDir = filepath.Join("/", ".pivot_root")
    if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
        return fmt.Errorf("Umount pivot_root dir %v", err)
    }

    return os.Remove(pivotDir)

}

func setUpMount() {
    //获取当前路径
    pwd, err := os.Getwd()
    if err != nil {
        log.Errorf("Get current location error %v", err)
        return
    }
    log.Infof("Current location is %s", pwd)
    pivotRoot(pwd)

    defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

    syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
    syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID | syscall.MS_STRICTATIME, "mode=755")
}
