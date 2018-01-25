package main

import (

    "os/exec"
    "syscall"
    "os"
    "log"
    "strconv"
    "io/ioutil"
    "fmt"
    "path"
)

const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
/*    fmt.Println("args : ", os.Args[0])
//    cmd := exec.Command("sh")
//    cmd := exec.Command("/proc/self/exe")
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
    }
//    cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(1), Gid: uint32(1)}
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err!= nil {
        log.Fatal(err)
    }
    os.Exit(-1)
*/
   
//   if os.Args[0] == "/proc/self/exe" {


//   }
    
    if os.Args[0] == "/proc/self/exe" {
       //容器进程
       fmt.Println("current pid : ", syscall.Getpid())
       fmt.Println()
       cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
       cmd.SysProcAttr = &syscall.SysProcAttr{}
       cmd.Stdin = os.Stdin
       cmd.Stdout = os.Stdout
       cmd.Stderr = os.Stderr
 
//       fmt.Println("cmd.Process.Pid : ", cmd.Process.Pid)

/*       if err := cmd.Run(); err != nil {
           fmt.Println(err)
           log.Fatal(err)
           os.Exit(1) 
       }
*/    
      if err := cmd.Start(); err != nil {
          fmt.Println(err)
          log.Fatal(err)
          os.Exit(1)
      }else {
          fmt.Println("cmd.Process.Pid : ", cmd.Process.Pid)
          
      }

  
   }
    
    cmd := exec.Command("/proc/self/exe")
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
    }
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Start(); err!= nil {
        fmt.Println("ERROR", err)
        os.Exit(1)
    }else {
        fmt.Println("external process id :", cmd.Process.Pid)
        
        //在系统默认的cgroup路径中(hierachy & root cgroup)，创建子cgroup  
        os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)
        //将容器进程加入到这个cgroup中
        ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
       //限制cgroup的进程使用
       ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644) 
   }
   cmd.Process.Wait()
}
