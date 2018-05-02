package main

import (
    log "github.com/Sirupsen/logrus"
    "fmt"
    "os"
    "io/ioutil"
//    "encoding/json"
//    "text/tabwriter"
//    "os/exec"
    "dockerEngine/src/container"
)

func logContainer(containerName string) {
    dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)  //string
    logFileLocation := dirUrl + container.ContainerLogFile

    file, err := os.Open(logFileLocation)
    defer file.Close()
    if err != nil {
        log.Errorf("Log container open file %s error %v", logFileLocation, err)
        return
    }

    content, err := ioutil.ReadAll(file)
    if err != nil {
        log.Errorf("Log container read file %s error %v", logFileLocation, err)
        return
    }

    fmt.Fprint(os.Stdout, string(content))
}
