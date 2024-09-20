package config

import (
    "encoding/json"
    "fmt"
    "os"
)

type ServerConf struct {
    BindAddr string `json:"bind_addr"`
    BindPort int    `json:"bind_port"`
    CertFile string `json:"cert_file"`
    KeyFile  string `json:"key_file"`
    CaCert   string `json:"ca_cert"`
    Token    string `json:"token"`
}

func LoadServerConfig(configPath string) (*ServerConf, error) {
    var conf ServerConf
    jsonStr, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read server config: %w", err)
    }
    if err := json.Unmarshal(jsonStr, &conf); err != nil {
        return nil, fmt.Errorf("failed to unmarshal server config: %w", err)
    }

    // Add basic validation
    if conf.BindAddr == "" || conf.BindPort == 0 {
        return nil, fmt.Errorf("invalid configuration: BindAddr or BindPort missing")
    }

    return &conf, nil
}