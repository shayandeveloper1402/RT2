package server

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "log"
    "net"
    "os"
    "pzrp/pkg/config"
    pkgErr "pzrp/pkg/errors"
    "sync"
)

// Adding resource cleanup and error logging in server.go

type TCPClientNode struct {
    *tcp.TCPNode
}

func (t *TCPClientNode) cleanup() {
    t.TCPNode.Close()
}

func NewTCPClientNode(conn tcp.DuplexConnection, readCtx, writeCtx context.Context) *TCPClientNode {
    node := &TCPClientNode{
        TCPNode: tcp.NewTCPNode(conn, readCtx, writeCtx, true),
    }
    node.Pack = node.overridePack
    node.UnPack = node.overrideUnPack
    return node
}

func Run(ctx context.Context, conf *config.ServerConfig) error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Server panic recovered: %v", r)
        }
    }()

    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{conf.ServerCert},
        RootCAs:      x509.NewCertPool(),
        MinVersion:   tls.VersionTLS13,
    }

    ln, err := tls.Listen("tcp", conf.ListenAddress, tlsConfig)
    if err != nil {
        return fmt.Errorf("failed to start server: %w", err)
    }
    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }

        node := NewTCPClientNode(conn, ctx, ctx)
        defer node.cleanup()
        // Actual server processing...
    }
    return nil
}