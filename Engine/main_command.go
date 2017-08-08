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
       resConf := &subsystems.ResourceConfig{
           MemoryLimit: context.String("m"),
           CpuSet: context.String("cpuset"),
           CpuShare: context.String("cpushare"),
       }
       volume := context.String("v")
 
       Run(tty, cmdArray, resConf, volume)
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
