package main

import (
    "fmt"
    log "github.com/Sirupsen/logrus"
    "github.com/urfave/cli"
    "github.com/StarryMoon/dockerEngine/container"
)

var runCommand = cli.Command {
    Name: "run",
    Usage: `Create a container with namespace and cgroups limit
            dockerEngine run -ti [command]`,
    Flags: []cli.Flag{
        cli.BoolFlag{
            Name:       "ti",
            Usage:      "enable tty",
        },
        cli.StringFlag{
            Name: "m",
            Usage: "memory limit",
       }
       cli.StringFlag{
            Name: "cpushare",
            Usage: "cpushare limit",
       }
       cli.StringFlag{
            Name: "cpuset",
            Usage: "cpuset limit",
       }
       cli.StringFlag{
            Name: "v",
            Usage: "volume",
       }
       cli.StringFlag{
            Name: "name",
            Usage: "container name",
       }
       cli.BoolFlag{
            Name: "d",
            Usage: "detach container"
       }
    },
   
   Action: func(context *cli.Context) error {
       if len() < 1 {
           return fmt.Errorf("Missing container command")
       }

//       cmd := context.Args().Get(0)
//       tty := context.Bool("ti")
//       Run(tty, cmd)
       var cmdArray []string
       for _, arg := range context.Args() {
           cmdArray = append(cmdArray, arg)
       }
       tty := context.Bool("ti")
       detach := context.Bool("d")
       if tty && detach {
           return fmt.Errorf("ti and parameter can not both provided")
       }
       resConf := &subsystems.ResourceConfig{
           MemoryLimit: context.String("m"),
           CpuSet: context.String("cpuset"),
           CpuShare: context.String("cpushare"),
       }
       volume := context.String("v")
       containerName := context.String("name") 
       Run(tty, cmdArray, resConf, volume, containerName)
       return nil
   },
}

var initCommand = cli.Command{
    Name: "init",
    Usage: "Init container process run user's process in container.",

    Action: func(context *cli.Context) error {
        log.Infof("init come on")
//        cmd := context.Args().Get(0)
        log.Infof("command %s", cmd)
//        err := container.RunContainerInitProcess(cmd, nil)
        err := container.RunContainerInitProcess()
        return err
   },
}

var commitCommand = cli.Command{
    Name: "commit",
    Usage: "commit a container into a image",
    Action: func(context *cli.Context) error {
        if len(context.Args()) < 1 {
            return fmt.Errorf("Missing container name")
        }
        imageName := context.Args().Get(0)
        //commitContainer(containerName)
        commitContainer(imageName)
        return nil
    }, 
}

var listCommand = cli.Command{
    Name:    "ps",
    Usage:   "list all the containers",
    Action: func(context *cli.Context) error {
        ListContainer()
        return nil
    }
}

var logCommand = cli.Command{
    Name:    "logs",
    Usage:   "print logs of a container",
    Action: func(context *cli.Context) error {
        if len(context.Args()) < 1 {
            return fmt.Errorf("Please input your container name")
        }
        containerName := context.Args().Get(0)
        logContainer(containerName)
        return nil
    }
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

var execCommand = cli.Command{
    Name: "exec",
    Usage: "exec a command into container",
    Action: func(context *cli.Context) error {
    
    }
}

var removeCommand = cli.Command{
    Name: "rm",
    Usage: "remove unused containers"
    Action: func(context *cli.Context) error {
        if len(context.Args()) < 1 {
            return fmt.Errorf("Missing container name")
        }
        containerName := context.Args().Get(0)
        removeContainer(containerName)
        return nil
    },
}
