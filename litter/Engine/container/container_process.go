package container

import (
    log "github.com/Sirupsen/logrus"
    "os"
    "syscall"
    "os/exec"
)

var (
    RUNNING               string = "running"
    STOP                  string = "stopped"
    EXIT                  string = "exited"
    DefaultInfoLocation   string = "/var/run/dockerEngine/%s/"
    ConfigName            string = "config.json"
    ContainerLogFile      string = "container.log"
)

type ContainerInfo struct {
    pid              string `json:"pid"`
    Id               string `json:"id"`
    Name             string `json:"name"`
    Command          string `json:"command"`
    CreateTime       string `json:"createTime"`
    Status           String `json:"status"`
}

func NewParentProcess(tty bool, containerName string) (*exec.Command, *os.File) {
    //args := []string{"init", command}
    readPipe, writePipe, err := NewPipe()
    if err != nil {
        log.Errorf("New pipe error %v", err)
        return nil, nil 
    }    
    cmd := exec.Command("/proc/self/exe", "init")
    cmd.SysProcAttr = &syscall.SysProcAttr{
        CloneFlags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS | syscall.CLONE_NET,
    }

    if tty {
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    }else {
        dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
        if err := os.MkdirAll(dirURL, 0622); err != nil {
            log.Errorf("NewParentProcess mkdir %s error %v", dirURL, err)
            return nil, nil
        }
        stdLogFilePath := dirURL + ContainerLogFile
        stdLogFile, err := os.Create(stdLogFilePath)
        if err != nil {
            log.Errorf("NewParentProcess create file %s error %v", stdLogFilePath, err)
            return nil, nil
        }
        cmd.Stdout = stdLogFile
    }
    cmd.ExtraFiles = []*os.File{readPipe}
//    cmd.Dir = "/root/busybox"
    mntURL := "/root/mnt/"
    rootURL := "/root/"
    NewWorkSpace(rootURL, mntURL)
    cmd.Dir = mntURL
    return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
    read, write, err := os.Pipe()
    if err != nil {
        return nil, nil, err
    }
    return read, write, nil
}


func NewWorkSpace(rootURL string, mntURL string, volume string) {
    CreateReadOnlyLayer(rootURL)
    CreateWriteLayer(rootURL)
    CreateMountPoint(rootURL, mntURL)
    if volume != "" {
        volumeURLs := volumeUrlExtract(volume)
        length := len(volumeURLs)
        if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
            MountVolume(rootURL, mntURL, volumeURLs)
            log.Infof("%q", volumeURLs)
        }else {
            log.Infof("volume parameter input is not correct.")
        }
    }
}

func CreateReadOnlyLayer(rootURL string) {
    busyboxURL := rootURL + "busybox/"
    busyboxTarURL := rootURL + "busybox.tar"
    exist, err := PathExists(busyboxURL)
    if err != nil {
        log.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
    }
    if exist == false {
        if err := os.Mkdir(busyboxURL, 0777); err != nil {
            log.Errorf("Mkdir dir %s error. %v", busyboxURL, err) 
        }
        if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
            log.Errorf("unTar dir %s error %v", busyboxTarURL, err)
        }
    }
}


func CreateWriteLayer(rootURL string) {
    writeURL := rootURL + "writeLayer/"
    if err := os.Mkdir(writeURL, 0777); err != nil {
        log.Errorf("Mkdir dir %s error. %v", writeURL, err)
    }
}

func MountVolume(rootURL string, mntURL string, volumeURLS {}string) {
    parentUrl := volumeURLs[0]
    if err := os.Mkdir(parentUrl, 0777); err != nil {
        log.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
    }
    containerUrl := volumeURLs[1]
    containerVolumeURL := mntURL + containerUrl
    if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
        log.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
    }
    dirs = "dirs=" + parentUrl
    cmd :+ exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr  
    if err := cmd.Run(); err != nil {
        log.Errorf("Mount volume failed. %v", err)
    }
}

func CreateMountPoint(rootURL string, mntURL string) {
    if err := os.Mkdir(mntURL, 0777); err != nil {
        log.Errorf("Mkdir dir %s error. %v", mntURL, err)
    }
    dirs := "dirs=" + rootURL + "writeLayer:" + "rootURL" + "busybox"
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


func DeleteWorkSpace(rootURL string, mntURL string, volume string) {
    if volume != "" {
        volumeURLs := volumeUrlExtract(volume)
        length := len(volumeURLs)
        if (length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "") {
            DeleteMountPointWithVolume(rootURL, mntURL, volumeURLs)
        }
    }

    DeleteMountPoint(rootURL, mntURL)
    DeleteWriteLayer(rootURL)
}

func DeleteMountPointWithVolume(rootURL string, mntURL string, volumeURLs []string) {
    containerUrl := mntURL + volumeURLs[1]
    cmd := exec.Command("umount", containerUrl)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("Umount volume failed. %v", err)
    }
}

func DeleteMountPoint(rootURL string, mntURL string) {
    cmd := exec.Command("umount", mntURL)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("%v", err)
    }
    if err := os.RemoveAll(mntURL); err != nil {
        log.Errorf("Remove dir %s error %v", mntURL, err)
    }
}

func DeleteWriteLayer(rootURL string) {
    writeURL := rootURL + "writeLayer/"
    if err := os.RemoveAll(writeURL); err != nil {
        log.Errorf("Remove dir %s error %v", writeURL, err)
    }
}
