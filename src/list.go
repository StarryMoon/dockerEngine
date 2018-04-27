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

    w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
    fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
    for _, item := range containers {
        fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
            item.Id,
            item.Name,
            item.Pid,
            item.Status,
            item.Command,
            item.CreateTime)
    }
    
    if err := w.Flush(); err != nil {
        log.Errorf("Flush error %v", err)
        return
    }
}


func getContainerInfo(file os.FileInfo) (*container.ContainerInfo, error) {
    
    containerName := file.Name()
    configFileDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
    configFileDir = configFileDir + container.ConfigName
    content, err := ioutil.ReadFile(configFileDir)
    if err != nil {
        log.Errorf("Read file %s error %v", configFileDir, err)
        return nil, err
    }

    var containerInfo containerInfo
    if err := json.Unmarshal(content, &containerInfo); err != nil {
        log.Errorf("Json unmarshal error %v", err)
        return nil, err
    }

    retrun &containerInfo, nil
}
