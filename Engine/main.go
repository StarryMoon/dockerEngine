package main

import (

    log "github.com/Sirupsen/logrus"
    "github.com/urfave/cli"
    "os"
)

const usage = `this is a simple container runtime implementation.`

func main() {

    app := cli.NewApp()
    app.Name = "dockerEngine"
    app.Usage = usage

    app.Commands = []cli.Command{
        initCommand,
        runCommand,
        commitCommand,
        listCommand,
        logCommand,
        stopCommand,
        execCommand,
        removeCommand,
    }

    app.Before = func(context *cli.Context) error {
        log.SetFormatter(&log.JSONFormatter())
        log.SetOutput(os.Stdout)
        return nil
   }

   if err := app.Run(os.Args); err!= nil {
       log.Fatal(err)
   }
}
