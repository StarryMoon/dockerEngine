package main

import (
    log "github.com/Sirupsen/logrus"
    "dockerEngine/src/container"
    "fmt"
    "os/exec"
)

func commitContainer(containerName string, imageName string) {
    mntURL := fmt.Sprintf(container.MntUrl, containerName)
    imageTar := container.RootUrl + "/" + imageName + ".tar"
    fmt.Printf("%s", imageTar)
    if _, err := exec.Command("tar", "-czf", imageTar, mntURL).CombinedOutput(); err != nil {
        log.Errorf("Tar folder %s error %v", mntURL, err)
    }
}
