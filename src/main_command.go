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
        cli.BoolFlag{
            Name: "d",
            Usage: "detach container",
        }, 
        cli.StringFlag{
            Name: "v",
            Usage: "volume",
        },
        cli.StringFlag{
            Name: "name",
            Usage: "container name",
        },
    },
    Action: func(context *cli.Context) error {
        if len(context.Args()) < 1 {
            return fmt.Errorf("Missing container command")
        }
//        cmd := context.Args().Get(0)
        fmt.Println("context args : ", context.Args())

        // user command
        var cmdArray []string
        for _, arg := range context.Args() {
            cmdArray = append(cmdArray, arg)
            fmt.Println("cmdArray : ", cmdArray)
        }
        fmt.Println("context args Get(0) : ", context.Args().Get(0))
        fmt.Println("context args : ", cmdArray)


        tty := context.Bool("ti")
        volume := context.String("v")
        detach := context.Bool("d")
        if tty && detach {
            fmt.Errorf("ti and d parameter can not be both used.")
        }

        containerName := context.String("name")

        memorylimit := context.String("m")
        cpuset := context.String("cpuset")
        cpushare := context.String("cpushare")
        fmt.Println("context args memory: ", memorylimit)
        fmt.Println("context args cpuset: ", cpuset)
        fmt.Println("context args cpushare: ", cpushare)
        resConf := &subsystems.ResourceConfig{
            MemoryLimit: memorylimit,
            CpuSet: cpuset,
            CpuShare: cpushare,
        }
//        RunCmd(tty, cmdArray, resConf)
//        Run(tty, cmdArray, volume)

          Run(tty, cmdArray, resConf, volume, containerName)
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
        err := container.RunContainerInitProcess()
        return err
    }, 
}

var commitCommand = cli.Command{
    Name: "commit",
    Usage: "Commit a container into one image",
    Action: func(context *cli.Context) error {
        log.Infof("commit come on")
        if len(context.Args())<1 {
            return fmt.Errorf("Missing container name")
        }

        //the default container is running in /root/mnt
        imageName := context.Args().Get(0)
        commitContainer(imageName)
        return nil
    },
}

var listCommand = cli.Command{
    Name: "ps",
    Usage: "list all the containers",
    Action: func(context *cli.Context) error {
        ListContainers()
        return nil
    },
}
