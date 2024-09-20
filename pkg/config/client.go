package config

import (
    "encoding/json"
    "fmt"
    "os"
    "pzrp/pkg/proto"
    "strconv"
    "strings"
)

type ClientConf struct {
    ServerAddr string `json:"server_addr"`
    ServerPort int    `json:"server_port"`
    CertFile   string `json:"cert_file"`
    KeyFile    string `json:"key_file"`
    CaCert     string `json:"ca_cert"`
    Services   map[string]ClientServiceConf
    Token      string `json:"token"`
}

type ClientServiceConf struct {
    Type       string `json:"type"`
    LocalIP    string `json:"local_ip"`
    LocalPort  int    `json:"local_port"`
    RemotePort int    `json:"remote_port"`
}

func LoadClientConfig(configPath string) (*ClientConf, error) {
    var conf ClientConf
    jsonStr, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read client config: %w", err)
    }
    if err := json.Unmarshal(jsonStr, &conf); err != nil {
        return nil, fmt.Errorf("failed to unmarshal client config: %w", err)
    }

    // Add basic validation
    if conf.ServerAddr == "" || conf.ServerPort == 0 {
        return nil, fmt.Errorf("invalid configuration: ServerAddr or ServerPort missing")
    }

    return &conf, nil
}