package main

imort (
    "github.com/StarryMoon/dockerEngine/container"
    "os"
    log "github.com/Sirupsen/logrus"
    "strings"
    "github.com/StarryMoon/dockerEngine/cgroups"
    "github.com/StarryMoon/dockerEngine/cgroups"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig, volume string, containerName string) error {
    //parent := container.NewParentProcess(tty, command)
    parent, writePipe := container.NewParentProcess(tty, volume)
    if parent == nil {
        log.Errorf("New parent process error")
        return
    }
    if err := parent.Start(); err != nil {
         log.Error(err)
    }

    containerName, err := recordContainerInfo(parent.Process.Pid, cmdArray, containerName)
    if err != nil {
        log.Errorf("Record container info error %v", err)
        return
    }

    cgroupManager := cgroups.NewCgroupManager("dockerEngine-cgroup")
    defer cgroupManager.Destroy()
    cgroupManager.Set(res)
    cgroupManager.Apply(parent.Process.Pid)

    sendInitCommand(comArray, writePipe)
//    parent.Wait()
    if tty {
        parent.Wait()
    }
    mntURL := "/root/mnt"
    rootURL := "/root"
    container.DeleteWorkSpace(rootURL, mntURL, volume)
    os.Exit(0)
//    os.Exit(-1)
}

func recordContaineInfo(containerPID int, commandArray []string, containerName string) (string, error) {
    id := randStringBytes(10)
    createTime := time.Now().Format("2006-01-02 15:05:04")
    command := strings.Join(commandArray, "")
    if containerName == "" {
        containerName = id
    }
    containerInfo := &containerContainerInfo{
        Id:          id,
        Pid:         strconv.Itoa(containerPID),
        Command:     command,
        CreateTime:  createTime,
        Status:      container.RUNNING,
        Name:        containerName,
    }

    jsonBytes, err := json.Marshal(containerInfo)
    if err != nil {
         log.Errorf("Record container info error %v", err)
        return "", err
    }
    jsonStr := string(jsonBytes)

    dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
    if err := os.MkdirAll(dirUrl, 0622); err != nil {
        log.Errorf("Mkdir error %s error %v", dirUrl, err)
        return "", err
    }
    fileName := dirUrl + "/" + container.ConfigName
    file, err := os.Create(fileName)
    defer file.Close()
    if err != nil {
        log.Errorf("Create file %s error %v", fileName, err)
        return "", err
    }
    if _, err := file.WriteString(jsonStr); err != nil {
        log.Errorf("File write string error %v", err)
        return "", err
    }

    return containerName, nil
}

func deleteContainerInfo(containerId string) {
    dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerId)
    if err := os.RemoveAll(dirURL); err != nil {
        log.Errorf("Remove dir %s error %v", dirURL, err)
    }
}

func sendInitCommand(comArray []string, writePipe *os.File) {
    command := strings.Join(comArray, " ")
    log.Infof("command all is %s", command)
    writePipe.WriteString(command)
    writePipe.Close()
}

func randStringBytes(n int) string {
    letterBytes := "1234567890"
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}
