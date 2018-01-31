package main

import (
    "fmt"
    log "github.com/Sirupsen/logrus"
    "github.com/urfave/cli"
    "dockerEngine/src/container"
    "dockerEngine/src/cgroups/subsystems"
//    "dockerEngine/src/cgroups"
)

var runCommand = cli.Command{
    Name: "run",
    Usage: "Create a container with namespace and cgroups limit dockerEngine run -ti [command]",
    Flags: []cli.Flag{
        cli.BoolFlag{
            Name: "ti",
            Usage: "enable tty",
        },
        cli.StringFlag{
            Name: "m",
            Usage: "memory limit",
        },
        cli.StringFlag{
            Name: "cpushare",
            Usage: "cpushare limit",
        },
        cli.StringFlag{
            Name: "cpuset",
            Usage: "cpuset limit",
        },
    },
    Action: func(context *cli.Context) error {
        if len(context.Args()) < 1 {
            return fmt.Errorf("Missing container command")
        }
        cmd := context.Args().Get(0)
//        var cmdArray []string
//        for _, arg := range context.Args() {
//            cmdArray = append(cmdArray, arg)
//        }
        fmt.Println("context args Get(0) : ", context.Args().Get(0))
//        fmt.Println("context args : ", cmdArray)
        tty := context.Bool("ti")
        resConf := &subsystems.ResourceConfig{
            MemoryLimit: context.String("m"),
            CpuSet: context.String("cpuset"),
            CpuShare: context.String("cpushare"),
        }
        RunCmd(tty, cmd, resConf)
        return nil
    },
}

var initCommand = cli.Command{
    Name: "init",
    Usage: "Init container process run user's process in container. Do not call it outside", 
    Action: func(context *cli.Context) error {
        log.Infof("init come on")
        cmd := context.Args().Get(0)
        fmt.Println("init command args : ", cmd)
        log.Infof("command %s", cmd)
        err := container.RunContainerInitProcess(cmd, nil)
        return err
    }, 
}
