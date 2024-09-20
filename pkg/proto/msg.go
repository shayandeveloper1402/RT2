package proto

import (
    "bytes"
    "encoding/binary"
    "errors"
)

// Define the PacketHead structure (adjust fields based on your requirements)
type PacketHead struct {
    Type   uint16 // Example field representing the type of packet
}

// Define the Packet structure
type Packet struct {
    Length uint16     // Length of the packet
    Head   PacketHead // Packet header
    Body   []byte     // Packet body (data)
}

var pkgOrder = binary.BigEndian
var pkgHeadOffset = 2 + binary.Size(PacketHead{})

// NewPacket function to parse the packet data
func NewPacket(data []byte) (*Packet, int, error) {
    if len(data) < pkgHeadOffset {
        return nil, 0, errors.New("incomplete data packet")
    }
    var (
        buf = bytes.NewBuffer(data)
        pkg = Packet{}
    )
    
    // Read the packet length
    if err := binary.Read(buf, pkgOrder, &pkg.Length); err != nil {
        return nil, 0, err
    }
    
    size := int(pkg.Length) + pkgHeadOffset
    if len(data) < size {
        return nil, 0, errors.New("data packet size mismatch")
    }
    
    // Read the packet header (adjust based on actual header size)
    if err := binary.Read(buf, pkgOrder, &pkg.Head); err != nil {
        return nil, 0, err
    }
    
    // The remaining data is the body of the packet
    pkg.Body = buf.Bytes()
    
    return &pkg, size, nil
}
