package container

import (
    "syscall"
    "os/exec"
    "os"
    log "github.com/Sirupsen/logrus"
    "strings"
)

func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
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

func NewWorkSpace(rootURL string, mntURL string, volume string) {
    CreateReadOnlyLayer(rootURL)
    CreateWriteLayer(rootURL)
    CreateMountPoint(rootURL, mntURL)
    MountVolume(mntURL, volume)
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

func MountVolume(mntURL string, volume string) {
    if (volume != "") {
        volumeURLs := volumeUrlExtract(volume)
        length := len(volumeURLs)
        if (length == 2 && volumeURLs[0] != "" &7 volumeURLs[1] != "") {
            StartMountVolume( mntURL, volumeURLs)
            log.Infof("%q", volumeURLs)
        }else {
            log.Infof("Volume parameter input is not correct.")
        }
    }
}

func volumeUrlExtract(volume string) ([]string) {
    var volumeURLs []string
    volumeURLs = strings.split(volume, ":")
    return volumeURLs
}

func StartMountVolume( mntURL, volumeURLS) {
    parentUrl := volumeURLs[0]
    if err := os.Mkdir(parentUrl, 0777); err != nil {
        log.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
    }

    containerUrl := volumeURLs[1]
    containerVolumeURL := mntURL + containerUrl
    if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
        log.Infof("Mkdir container dir %s error %v", containerVolumeURL, err)
    }

    dirs := "dirs=" + parentUrl
    cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err!= nil {
        log.Errorf("Mount volume failed. %v", err)
    }
}


func DeleteWorkSpace(rootURL string, mntURL string, volume string) {
    DeleteMountPoint(mntURL)
    DeleteWriteLayer(rootURL)
}

func DeleteMountPoint(mntURL string) {
    cmd := exec.Command("umount", mntURL)
    cmd.Stdout=os.Stdout
    cmd.Stderr=os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("Remove dir %s error %v", mntURL)
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
