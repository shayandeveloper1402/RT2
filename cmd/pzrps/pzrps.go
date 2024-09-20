package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "pzrp/pkg/config"
    "pzrp/server"
)

var VERSION string

func main() {
    configPath := flag.String("config", "pzrps.json", "configuration file path")
    showVersion := flag.Bool("version", false, "display version number")
    flag.Parse()

    if *showVersion {
        fmt.Println(VERSION)
        return
    }

    conf, err := config.LoadServerConfig(*configPath)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    if err := server.Run(ctx, conf); err != nil {
        log.Fatalf("Server encountered an error: %v", err)
    }
}