package main

import (
    "dockerEngine/src/container"
    "dockerEngine/src/cgroups"
    "dockerEngine/src/cgroups/subsystems"
    "os"
    "strings"
    log "github.com/Sirupsen/logrus"
    "math/rand"
    "time"
    "fmt"
    "encoding/json"
    "strconv"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume string, containerName string, imageName string) {
    id := randContainerIdGenerator()
    if containerName == "" {
        containerName = id
    }

    parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName)
    if parent == nil {
        log.Errorf("New parent process error")
        return
    }
    if err := parent.Start(); err != nil {
        log.Error(err)
    }
   
    containerName, err := recordContainerInfo(parent.Process.Pid, comArray, containerName, id, imageName, volume)
    if err != nil {
        log.Errorf("Record container info error %v", err)
        return
    }

    /* Do support cgroup arguments in cmd line
     *
    */

    // create cgroup
    cgroupManager := cgroups.NewCgroupManager("dockerEngine-cgroup")
    defer cgroupManager.Destory()
    cgroupManager.Set(res)
    cgroupManager.Apply(parent.Process.Pid)

    sendInitCommand(comArray, writePipe)

    // the parent process will wait the child process
    if tty {
        parent.Wait()
        deleteContainerInfo(containerName)
//        rootURL := "/root/"
//        mntURL := "/root/mnt/"
//        container.DeleteWorkSpace(rootURL, mntURL, volume, imageName)
        container.DeleteWorkSpace(volume, containerName)
        os.Exit(0)
    }
}

/*
func RunCmd(tty bool, command []string, res *subsystems.ResourceConfig) {
    parent, writePipe := container.NewParentProcess(tty)
    if parent == nil {
        log.Errorf("New parent process error")
        return
    }
    if err := parent.Start(); err != nil {
        log.Error(err)
    }

    // create cgroup
    cgroupManager := cgroups.NewCgroupManager("dockerEngine-cgroup")
    defer cgroupManager.Destory()
    cgroupManager.Set(res)
    cgroupManager.Apply(parent.Process.Pid)

    sendInitCommand(command, writePipe)
    parent.Wait()
    os.Exit(-1)
}
*/

func sendInitCommand(cmdArray []string, writePipe *os.File) {
    command := strings.Join(cmdArray, " ")
    log.Infof("command all is %s ", command)
    writePipe.WriteString(command)
    writePipe.Close()
}

func recordContainerInfo(containerPID int, commandArray []string, containerName string, containerId string, imageName string, volume string) (string, error) {
   // id := randContainerIdGenerator()
    createTime := time.Now().Format("2006-01-02 15:04:05")
    command := strings.Join(commandArray, "")
   // if containerName == "" {
   //     containerName = id
   // }

    containerInfo := &container.ContainerInfo{
        Id:              containerId,                 //id,
        Pid:             strconv.Itoa(containerPID),    //类型转换  int-->string
        Command:         command,
        CreateTime:      createTime,
        Status:          container.RUNNING,
        Name:            containerName,
        Image:           imageName,
        Volume:          volume,
    }

    jsonBytes, err := json.Marshal(containerInfo)
    if err != nil {
        log.Errorf("Json container info error %v", err)
        return "", err
    }

    //转成string
    jsonStr := string(jsonBytes)

    //文件路径
    dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
    if _, err := os.Stat(dirUrl); err != nil {
        if os.IsNotExist(err) {
            if fileErr := os.MkdirAll(dirUrl, 0622); fileErr != nil {
                log.Errorf("Mkdir error %s error %v", dirUrl, err)
                return "", fileErr
            }
        }else {
            return "", err
        }
    }

    fileName := dirUrl + "/" + container.ConfigName
    file, err := os.Create(fileName)
    defer file.Close()
    if err != nil {
        log.Errorf("Create config file %s error %v", fileName, err)
        return "", err
    }

    if _, err := file.WriteString(jsonStr); err != nil {
        log.Errorf("File write string error %v", err)
        return "", err
    }

    
    return containerName, nil
}

func deleteContainerInfo(containerName string) {
    dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
    if err := os.RemoveAll(dirUrl); err != nil {
        log.Errorf("Remove dir %s error %v", dirUrl, err)
    }
}

func randContainerIdGenerator() string {
    letterBytes := "1234567890"
    rand.Seed(time.Now().UnixNano())   //随机数种子
    b := make([]byte, 10)              //id 默认只有10位
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }

    return string(b)
}
