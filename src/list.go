package main

import (
    log "github.com/Sirupsen/logrus"
    "fmt"
    "os/exec"
    "dockerEngine/src/container"
)

func listContainer() {
    dirUrl := fmt.Sprintf(container.DefaultInfoLocation, "")  //string
    fmt.Println("list container info path : %s", dirUrl)
    dirUrl := dirUrl[:len(dirUrl)]   //???

    files, err := ioutil.ReadDir(dirUrl)    // []os.FileInfo
    if err != nil {
        log.Errorf("Read dir %s error %v", dirUrl, err)
        return
    }

    var containers []*container.ContainerInfo
    for _, file := range files {
        tmpContainer, err := getContainerInfo(file)
        if err != nil {
            log.Errorf("Get container info error %v", err)
            continue
        }
        containers =append(containers, tmpContainer)
    }
   
    

}


func getContainerInfo(file os.FileInfo) (*container.ContainerInfo, error) {
    
    containerName := file.Name()
    configFileDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
    configFileDir = configFileDir + container.ConfigName
}
