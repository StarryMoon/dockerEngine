package subsystems

import (
    "fmt"
    "os"
    "io/ioutil"
    "path"
    "strconv"
)

type CpuSetSubsystem struct {
}

func (s *CpuSetSubsystem) Set(cgroupPath string, res *ResourceConfig) error {
    if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
        if res.CpuSet != "" {
            if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
                return fmt.Errorf("set cgroup cpuset fail %v", err)
            }
        }
        return nil
    }else {
        return err
    }
}

func (s *CpuSetSubsystem) Remove(cgroupPath string) error {
    if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
        return os.RemoveAll(subsysCgroupPath)
    }else {
         return err
     }
}

func (s *CpuSetSubsystem) Apply(cgroupPath string, pid int) error {
    if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
        if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(strconv.Itoa(0)), 0644); err != nil {
            return fmt.Errorf("set cgroup cpuset cpus fail %v", err)
        }
        if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.mems"), []byte(strconv.Itoa(0)), 0644); err != nil {
            return fmt.Errorf("set cgroup cpuset mems fail %v", err)
        }
        if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
            return fmt.Errorf("set cgroup proc fail %v", err)
        }
        return nil
    }else {
        return fmt.Errorf("get cpuset cgroup %s error: %v", cgroupPath, err)
    }
}

func (s *CpuSetSubsystem) Name() string {
    return "cpuset"
}
