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

    //cmd.Dir = "/root/busybox"
    
    //specify the foot file system
    mntURL := "/root/mnt/"
    rootURL := "/root/"
    NewWorkSpace(rootURL, mntURL)
    cmd.Dir = mntURL

    return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
    read, write, err := os.Pipe()
    if err != nil {
        return nil, nil, nil 
    }
    return read, write, nil
}

func NewWorkSpace(rootURL string, mntURL string) {
    CreateReadOnlyLayer(rootURL)
    CreateWriteLayer(rootURL)
    CreateMountPoint(rootURL, mntURL)
}


//decompress the tar package of read only file system
func CreateReadOnlyLayer(rootURL string) {
    busyboxURL := rootURL + "busybox/"
    busyboxTarURL := rootURL + "busybox.tar"
    exist, err := PathExists(busyboxURL)
    if err != nil {
        log.Infof("Fail to judge whether the directory %s exists. %v ", busyboxURL, err)
    }
    
    if exist == false {
        if err := os.Mkdir(busyboxURL, 0777); err != nil {
            log.Errorf("Mkdir directory %s error.  %v", busyboxURL, err)
        }

        if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
            log.Errorf("unTar directory %s error %v", busyboxTarURL, err)
        }
    }

}

//create a read-write layer
func CreateWriteLayer(rootURL string) {
    writeURL := rootURL + "writeLayer/"
    if err := os.Mkdir(writeURL, 0777); err != nil {
        log.Errorf("Mkdir directory %s error %v", writeURL, err)
    }
}

//mount rw layer and ro layer
func CreateMountPoint(rootURL string, mntURL string) {
    //create mnt mount point
    if err := os.Mkdir(mntURL, 0777); err != nil {
        log.Errorf("Mkdir directory %s error %v", mntURL, err)
    }

    //"mount -t aufs -o dirs=/root/writeLayer:/root/busybox/ none /root/mnt"
    dirs := "dirs=" + rootURL + "writeLayer:" + rootURL + "busybox"
    cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("%v", err)
    }
}

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }

    if os.IsNotExist(err) {
        return false, nil
    }

    return false, err
}
