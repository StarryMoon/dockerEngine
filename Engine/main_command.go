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
            Name:       "ti"
            Usage:      "enable tty"
        },
    },
   
   Action: func(context *cli.Context) error {
       if len() < 1 {
           return fmt.Errorf("Missing container command")
       }

       cmd := context.Args().Get(0)
       tty := context.Bool("ti")
       Run(tty, cmd)
       return nil
   },
}

var initCommand = cli.Command{
    Name: "init",
    Usage: "Init container process run user's process in container.",

    Action: func(context *cli.Context) error {
        log.Infof("init come on")
        cmd := context.Args().Get(0)
        log.Infof("command %s", cmd)
        err := container.RunContainerInitProcess(cmd, nil)
        return err
   },
}
