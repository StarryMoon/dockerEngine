package main

import (
    "fmt"
    log "github.com/Sirupsen/logrus"
    "github.com/urfave/cli"
    "dockerEngine/src/container"
    "dockerEngine/src/cgroups/subsystems"
    "os"
)

var runCommand = cli.Command{
    Name: "run",
    Usage: "Create a container with namespace and cgroups limit dockerEngine run -ti [command]",
    Flags: []cli.Flag{
        cli.BoolFlag{
            Name: "t",
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
        cli.StringFlag{
            Name: "i",
            Usage: "image name",
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


        tty := context.Bool("t")
        volume := context.String("v")
        detach := context.Bool("d")
        if tty && detach {
            fmt.Errorf("ti and d parameter can not be both used.")
        }

        containerName := context.String("name")
        imageName := context.String("i")

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

          Run(tty, cmdArray, resConf, volume, containerName, imageName)
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
        if len(context.Args())<2 {
            return fmt.Errorf("Missing container name")
        }

        //the default container is running in /root/mnt
        containerName := context.Args().Get(0)
        imageName := context.Args().Get(1)
        commitContainer(containerName, imageName)
        return nil
    },
}

var listCommand = cli.Command{
    Name: "ps",
    Usage: "list all the containers",
    Action: func(context *cli.Context) error {
        listContainers()
        return nil
    },
}

var logCommand = cli.Command{
    Name: "log",
    Usage: "print logs of a container",
    Action: func(context *cli.Context) error {
        if len(context.Args()) <1 {
            return fmt.Errorf("Please input your container name")
        }
        containerName := context.Args().Get(0)
        logContainer(containerName)
        return nil
    },
}


// first : to analyze the containerName and cmd
// second: to find out the correspondingly pid
// third : to set up envirnment varibles
// fourth: to execute the "exec" again
// fifth : to trigger the "setns"
// ***note: environment varibles is a threshold, if it is set, the exec will not be executed again; and the content of "exec" is really executed by the "setns"
var execCommand = cli.Command{
    Name: "exec",
    Usage: "exec a command into container",
    Action: func(context *cli.Context) error {
        if os.Getenv(ENV_EXEC_CMD) != "" {    //if true, execute the cmd directly, not to create the environment again
            log.Infof("pid callback pid %s", os.Getpid())
            return nil
        }

        if len(context.Args()) < 2 {
            return fmt.Errorf("Missing container name or command")
        }

        containerName := context.Args().Get(0)
        var commandArray []string
        for _, arg := range context.Args().Tail() { //except the first element
            commandArray = append(commandArray, arg)
        }

        execContainer(containerName, commandArray)
        return nil
    },
}

var stopCommand = cli.Command{
    Name: "stop",
    Usage: "stop a container",
    Action: func(context *cli.Context) error {
        if len(context.Args()) < 1 {
            return fmt.Errorf("Missing container name")
        }
        containerName := context.Args().Get(0)
        stopContainer(containerName)
        return nil
    },
}

var deleteCommand = cli.Command{
    Name: "rm",
    Usage: "remove specified container",
    Action: func(context *cli.Context) error {
        if len(context.Args()) < 1 {
            return fmt.Errorf("Missing container name")
        }

        containerName := context.Args().Get(0)
        deleteContainer(containerName)
        return nil
    },
}
