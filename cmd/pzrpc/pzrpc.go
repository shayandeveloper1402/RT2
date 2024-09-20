package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "pzrp/client"
    "pzrp/pkg/config"
)

var VERSION string

func main() {
    configPath := flag.String("config", "pzrpc.json", "configuration file path")
    showVersion := flag.Bool("version", false, "display version number")
    flag.Parse()

    if *showVersion {
        fmt.Println(VERSION)
        return
    }

    conf, err := config.LoadClientConfig(*configPath)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    if err := client.Run(ctx, conf); err != nil {
        log.Fatalf("Client encountered an error: %v", err)
    }
}