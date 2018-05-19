package main

import (
    log "github.com/Sirupsen/logrus"
    "fmt"
    "os"
    "dockerEngine/src/container"
)

func deleteContainer(containerName string) {
    containerInfo, err := getContainerInfoByName(containerName)
    if err != nil {
        log.Errorf("Delete container getContainerInfoByName %s error %v", containerName, err)
        return
    }

    if containerInfo.Status != container.STOP {
        log.Errorf("Couldn't remove running container")
        return
    }
    
    dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)  // runtime info
    if err := os.RemoveAll(dirURL); err != nil {
        log.Errorf("Remove file %s error", dirURL, err)
        return
    }
    container.DeleteWorkSpace(containerInfo.Volume, containerName)      // image construct info
}

