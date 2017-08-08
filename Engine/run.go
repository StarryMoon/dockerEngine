package main

imort (
    "github.com/StarryMoon/dockerEngine/container"
    "os"
    log "github.com/Sirupsen/logrus"
    "strings"
    "github.com/StarryMoon/dockerEngine/cgroups"
    "github.com/StarryMoon/dockerEngine/cgroups"
)

func Run(tty bool, command string, res *subsystems.ResourceConfig, volume string) error {
    //parent := container.NewParentProcess(tty, command)
    parent, writePipe := container.NewParentProcess(tty, volume)
    if parent == nil {
        log.Errorf("New parent process error")
        return
    }
    if err := parent.Start(); err != nil {
         log.Error(err)
    }
    cgroupManager := cgroups.NewCgroupManager("dockerEngine-cgroup")
    defer cgroupManager.Destroy()
    cgroupManager.Set(res)
    cgroupManager.Apply(parent.Process.Pid)

    sendInitCommand(comArray, writePipe)
    parent.Wait()
    mntURL := "/root/mnt"
    rootURL := "/root"
    container.DeleteWorkSpace(rootURL, mntURL, volume)
    os.Exit(0)
//    os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
    command := strings.Join(comArray, " ")
    log.Infof("command all is %s", command)
    writePipe.WriteString(command)
    writePipe.Close()
}
