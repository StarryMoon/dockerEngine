package cgroups

import (
    "dockerEngine/src/cgroups/subsystems"
    "github.com/Sirupsen/logrus"
)

type CgroupManager struct {
    Path  string   //cgroup name
    Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
    return &CgroupManager{
        Path: path,
    }
}

func (c * CgroupManager) Apply(pid int) error {
    for _, subSysIns := range(subsystems.SubSystemsIns) {
        if err := subSysIns.Apply(c.Path, pid); err != nil {
            logrus.Errorf("apply cgroup fail %v", err)
        }    
    }
    return nil
}

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
    for _, subSysIns := range(subsystems.SubSystemsIns) {
        if err := subSysIns.Set(c.Path, res); err != nil {
            logrus.Errorf("set cgroup fail %v", err)
        } 
    }
    return nil
}

func (c *CgroupManager) Destory() error {
    for _, subSysIns := range(subsystems.SubSystemsIns) {
        if err := subSysIns.Remove(c.Path); err != nil {
            logrus.Errorf("remove cgroup fail %v", err)
        }
    }
    return nil
}
