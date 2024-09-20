package client

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "errors"
    "fmt"
    "log"
    "net"
    "os"
    "pzrp/pkg/config"
    pkgErr "pzrp/pkg/errors"
    "pzrp/pkg/proto"
    "sync"
)

// Adding resource cleanup and error logging in client.go

type clientNodeInfo struct {
    node        proto.Node
    readCancel  context.CancelFunc
    writeCancel context.CancelFunc
    readCtx     context.Context
    writeCtx    context.Context
}

type TunnelClientNode struct {
    *tcp.TCPNode
    serviceMapping sync.Map
}

func (t *TunnelClientNode) cleanup() {
    t.readCancel()
    t.writeCancel()
    t.node.Close()
}

func Run(ctx context.Context, conf *config.ClientConfig) error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Client panic recovered: %v", r)
        }
    }()

    // Set up TLS configuration for reuse
    tlsConfig := &tls.Config{
        RootCAs:            x509.NewCertPool(),
        InsecureSkipVerify: conf.InsecureSkipVerify,
        MinVersion:         tls.VersionTLS13,
    }

    // Example of reusing a TLS connection for performance
    conn, err := tls.Dial("tcp", conf.ServerAddress, tlsConfig)
    if err != nil {
        return fmt.Errorf("failed to connect: %w", err)
    }
    defer conn.Close()

    node := &TunnelClientNode{TCPNode: tcp.NewTCPNode(conn, ctx)}
    defer node.cleanup()

    // Actual implementation continues...
    return nil
}