package main

import (


)

func removeContainer(containerName strin) {
    containerInfo, err := getContainerInfoByName(containerName)
    if err != nil {
        log.Errorf("Get container %s info error %v", containerName, err)
        return
    }
    if containerInfo.Status != container.STOP {
        log.Errorf("must remove stop container")
        return
    }
    dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
    if err := os.RemoveAll(dirURL); err != nil {
        log.Errorf("Remove file %s error %v", dirURL, err)
        return
    }
}
