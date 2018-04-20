package main

import (
    log "github.com/Sirupsen/logrus"
    "github.com/urfave/cli"
    "os"
    "fmt"
)

const usage = `dockerEngine is a simple container runtime implementation. The purpose of this project is to learn how docker works and how to write a docker by ourselves. Enjoy it & just for fun.`

func main() {

    app := cli.NewApp()
    app.Name = "dockerEngine"
    app.Usage = usage

    //the command of cmd line
    app.Commands = []cli.Command{
        initCommand,  
        runCommand,
        commitCommand,
    }

    app.Before = func(context *cli.Context) error {
        log.SetFormatter(&log.JSONFormatter{})
        log.SetOutput(os.Stdout)
        return nil
    }

    osArgs := os.Args
    fmt.Println("os Args : ", osArgs) 
    if err := app.Run(osArgs); err!=nil {
        log.Fatal(err)
    }
}
