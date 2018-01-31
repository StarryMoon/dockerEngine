package subsystems

type ResourceConfig struct {
    MemoryLimit string
    CpuSet      string
    CpuShare    string
}

type Subsystem interface {
    Name() string
    Set(path string, res *ResourceConfig) error
    Apply(path string, pid int) error
    Remove(path string) error
}

//Global Variable
//SubSystem[0] = &cpuSet
//cpuSet = CpuSetSubsystem{}
var (
    SubSystemsIns = []Subsystem{
        &CpuSetSubsystem{},
        &CpuSubsystem{},
        &MemorySubsystem{},
    }
)
