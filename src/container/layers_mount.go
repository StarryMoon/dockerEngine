package container

import (
//    "syscall"
    "os/exec"
    "os"
    log "github.com/Sirupsen/logrus"
    "strings"
    "fmt"
)

/* ReadOnlyLayer: RootUrl + imageName
 * WriteLayer: WriteLayerUrl + containerName
 * MountLayer: MntUrl + containerName + containerURL
 */

func NewWorkSpace(volume string, containerName string, imageName string) {
    CreateReadOnlyLayer(imageName)
    CreateWriteLayer(containerName)
    StartMountLayers(containerName, imageName)
    MountVolume(volume, containerName)
}


//decompress the tar package of read only file system
func CreateReadOnlyLayer(imageName string) {
    unTarFoldUrl := RootUrl + "/" + imageName + "/"
    imageTarURL := RootUrl + "/" + imageName + ".tar"
    exist, err := PathExists(unTarFoldUrl)
    if err != nil {
        log.Infof("Fail to judge whether the directory %s exists. %v ", unTarFoldUrl, err)
    }
    
    if exist == false {
        if err := os.Mkdir(unTarFoldUrl, 0777); err != nil {
            log.Errorf("Mkdir directory %s error.  %v", unTarFoldUrl, err)
        }

        if _, err := exec.Command("tar", "-xvf", imageTarURL, "-C", unTarFoldUrl).CombinedOutput(); err != nil {
            log.Errorf("unTar directory %s error %v", imageTarURL, err)
        }
    }
}

//create a read-write layer
func CreateWriteLayer(containerName string) {
    writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
    if err := os.MkdirAll(writeURL, 0777); err != nil {                  //multiple file dirs need to be created
        log.Errorf("Mkdir directory %s error %v", writeURL, err)
    }
}

//mount rw layer and ro layer
func StartMountLayers(containerName string, imageName string) {
    //create mnt mount point
    mntURL := fmt.Sprintf(MntUrl, containerName)
    if err := os.MkdirAll(mntURL, 0777); err != nil {                    //multiple file dirs need to be created
        log.Errorf("Mkdir directory %s error %v", mntURL, err)
    }

    //"mount -t aufs -o dirs=/root/writeLayer:/root/busybox/ none /root/mnt"
    tmpWriteLayerUrl := fmt.Sprintf(WriteLayerUrl, containerName)
    tmpReadOnlyLayerUrl := RootUrl + "/" + imageName
    dirs := "dirs=" + tmpWriteLayerUrl + ":" + tmpReadOnlyLayerUrl

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

func MountVolume(volume string, containerName string) {
    if (volume != "") {
        volumeURLs := volumeUrlExtract(volume)
        length := len(volumeURLs)
        if (length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "") {
            mntURL := fmt.Sprintf(MntUrl, containerName)
            StartMountVolume(mntURL, volumeURLs)
            log.Infof("%q", volumeURLs)
        }else {
            log.Infof("Volume parameter input is not correct.")
        }
    }
}

func volumeUrlExtract(volume string) ([]string) {
    var volumeURLs []string
    volumeURLs = strings.Split(volume, ":")
    return volumeURLs
}

func StartMountVolume(mntURL string, volumeURLs []string) {
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


//func DeleteWorkSpace(rootURL string, mntURL string, volume string) {
func DeleteWorkSpace(volume string, containerName string) {
    if (volume != "") {
        volumeURLs := volumeUrlExtract(volume)
        mntURL := fmt.Sprintf(MntUrl, containerName)
        length := len(volumeURLs)
        if (length == 2 && volumeURLs[0]!= "" && volumeURLs[1]!= "") {
            DeleteVolume(mntURL, volumeURLs)
        }
    }
    DeleteMountPoint(containerName)
    DeleteWriteLayer(containerName)
}

func DeleteVolume(mntURL string, volumeURLs []string) {
    containerUrl := mntURL + "/" + volumeURLs[1]
    cmd := exec.Command("umount", containerUrl)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("Umount volume failed. %v", err)
    }
}

func DeleteMountPoint(containerName string) {
    mntURL := fmt.Sprintf(MntUrl, containerName)
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

func DeleteWriteLayer(containerName string) {
    writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
    if err := os.RemoveAll(writeURL); err != nil {
        log.Errorf("Remove dir %s error %v", writeURL, err)
    }
}
