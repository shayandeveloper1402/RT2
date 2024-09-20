package server

import (
    "context"
    "crypto/rand"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "pzrp/pkg/config"
    pkgErr "pzrp/pkg/errors"
    "pzrp/pkg/proto"
    "pzrp/pkg/proto/tcp"
    "pzrp/pkg/proto/udp"
    "pzrp/pkg/utils"
    "time"
)

type TunnelNode struct {
    *tcp.TCPNode
    services map[uint8]map[uint16]proto.Node
    ctx      context.Context
    cancel   context.CancelFunc
    key      []byte
}

func NewTunnelNode(conn tcp.DuplexConnection, ctx context.Context, key []byte) *TunnelNode {
    _ctx, _cancel := context.WithCancel(ctx)
    node := &TunnelNode{
        TCPNode:  tcp.NewTCPNode(conn, _ctx, _ctx, false),
        services: map[uint8]map[uint16]proto.Node{},
        ctx:      _ctx,
        cancel:   _cancel,
        key:      key,
    }
    node.Pack = node.overridePack
    node.UnPack = node.overrideUnPack
    return node
}

func (node *TunnelNode) AddServer(protocol uint8, port uint16) {
    services, ok := node.services[protocol]
    if !ok {
        services = make(map[uint16]proto.Node)
        node.services[protocol] = services
    }
    services[port] = nil
}

func (node *TunnelNode) findServer(protocol uint8, serverPort uint16) proto.Node {
    s1, ok := node.services[protocol]
    if !ok {
        return nil
    }
    s2, ok := s1[serverPort]
    if !ok {
        return nil
    }
    return s2
}

// Optimized to log errors instead of panic and ensure resource cleanup
func (node *TunnelNode) dispatchMsg() {
    logger := utils.GetLogger(node.ctx)
    defer func() {
        if r := recover(); r != nil {
            logger.Error("dispatch message failed", "error", r)
        }
        node.cancel()
    }()
    for {
        msg, err := node.Read()
        if err != nil {
            logger.Error("failed to read message", "error", err)
            return
        }
        nextNode := node.findServer(msg.Protocol, msg.ServerPort)
        if nextNode == nil {
            logger.Error("no server found for message", "protocol", msg.Protocol, "server_port", msg.ServerPort)
            continue
        }
        if err := nextNode.Write(msg); err != nil {
            logger.Error("failed to dispatch message", "error", err, "node", nextNode)
        }
    }
}

func (node *TunnelNode) collectMsg(server proto.Node) {
    logger := utils.GetLogger(node.ctx)
    defer func() {
        if r := recover(); r != nil {
            logger.Error("collect message failed", "error", r)
        }
    }()
    // Collect messages from server
    for {
        msg, err := server.Read()
        if err != nil {
            logger.Error("failed to read message from server", "error", err)
            return
        }
        if err := node.Write(msg); err != nil {
            logger.Error("failed to write message", "error", err)
        }
    }
}

func (node *TunnelNode) Run() {
    logger := utils.GetLogger(node.ctx)
    defer func() {
        if r := recover(); r != nil {
            logger.Error("tunnel abnormal exit", "error", r)
        }
        node.cancel()
    }()
    go node.initServer() // Running the server initialization concurrently
    node.TCPNode.Run()   // Starting the TCP node run loop
}

func (node *TunnelNode) startTCPServer(port uint16) proto.Node {
    addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
    if err != nil {
        log.Fatalf("failed to resolve TCP address: %v", err)
    }
    lis, err := net.ListenTCP("tcp4", addr)
    if err != nil {
        log.Fatalf("failed to start TCP server: %v", err)
    }
    defer lis.Close() // Ensure the listener is closed properly
    logger := utils.GetLogger(node.ctx).With("protocol", "tcp", "server_port", port)
    ctx := utils.SetLogger(node.ctx, logger)
    srv := NewTCPServerNode(lis, ctx)
    go srv.Run()
    return srv
}

func (node *TunnelNode) startUDPServer(port uint16) proto.Node {
    server, err := net.ListenUDP("udp", &net.UDPAddr{
        IP:   net.IPv4(0, 0, 0, 0),
        Port: int(port),
    })
    if err != nil {
        log.Fatalf("failed to start UDP server: %v", err)
    }
    defer server.Close() // Ensure the UDP server is closed properly
    srv := udp.NewUdpServerNode(server, port, node.ctx)
    go srv.Run()
    return srv
}