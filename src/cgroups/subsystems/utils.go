package subsystems

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "path"
)


// "/sys/fs/cgroup"
func FindCgroupMountPoint(subsystem string) string{
    f, err := os.Open("/proc/self/mountinfo")
    if err != nil {
        return ""
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        txt := scanner.Text()
        fields := strings.Split(txt, " ")   //To find the last array
        for _, opt := range strings.Split(fields[len(fields)-1], ",") {    //the last one, and then split in "," 
            if opt == subsystem {
                return fields[4]
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return ""
    }

    return ""
}


// "/sys/fs/cgroup/cpu"
// if the subsystem is not exist, auto to create it.
//
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) { 
    cgroupRoot := FindCgroupMountPoint(subsystem)
    fmt.Println("cgroup mount point : ", cgroupRoot)
    fmt.Println("cgroup mount point : ", path.Join(cgroupRoot, cgroupPath))
    if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
        if os.IsNotExist(err) {
            if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
            
            }else {
                 return "", fmt.Errorf("error create cgroup %v", err)
             }   
        }
        return path.Join(cgroupRoot, cgroupPath), nil
    }else {
         return "", fmt.Errorf("cgroup path error %v", err)
     }
}
