package config

// ServerConfig holds server configuration
type ServerConfig struct {
    Port int
}

// ClientConfig holds client configuration
type ClientConfig struct {
    ServerAddr string
    Username   string
}