package main

import (
    "flag"
    "fmt"
    "os"

    "vibhordubey333/road.sh/broadcast-server/internal/app"
    "vibhordubey333/road.sh/broadcast-server/internal/config"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: broadcast-server [start|connect]")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "start":
        startServer()
    case "connect":
        connectClient()
    default:
        fmt.Println("Unknown command. Use 'start' or 'connect'")
        os.Exit(1)
    }
}

func startServer() {
    port := flag.Int("port", 8080, "Port to listen on")
    flag.CommandLine.Parse(os.Args[2:])
    
    cfg := config.ServerConfig{
        Port: *port,
    }
    
    server := app.NewServer(cfg)
    server.Start()
}

func connectClient() {
    server := flag.String("server", "localhost:8080", "Server address to connect to")
    name := flag.String("name", "anonymous", "Your display name")
    flag.CommandLine.Parse(os.Args[2:])
    
    cfg := config.ClientConfig{
        ServerAddr: *server,
        Username:   *name,
    }
    
    client := app.NewClient(cfg)
    client.Connect()
}