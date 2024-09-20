package utils

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "net"
    "pzrp/pkg/proto"
    "time"
)

func TCPipe(port1 int, port2 int, delay time.Duration) (*net.TCPConn, *net.TCPConn) {
    addr1 := fmt.Sprintf("127.0.0.1:%d", port1)
    addr2 := fmt.Sprintf("127.0.0.1:%d", port2)
    listener, err := net.Listen("tcp", addr1)
    if err != nil {
        slog.Error("Failed to start TCP listener", "error", err, "address", addr1)
        return nil, nil
    }
    defer listener.Close()

    ch := make(chan net.Conn)
    go func(ch chan<- net.Conn) {
        time.Sleep(delay)
        conn2, err := net.Dial("tcp", addr2)
        if err != nil {
            slog.Error("Failed to connect", "error", err, "address", addr2)
            return
        }
        ch <- conn2
    }(ch)

    conn1, err := listener.Accept()
    if err != nil {
        slog.Error("Failed to accept connection", "error", err)
        return nil, nil
    }

    conn2 := <-ch
    return conn1.(*net.TCPConn), conn2.(*net.TCPConn)
}