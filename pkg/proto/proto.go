package proto


// Protocol constants
const (
    PROTO_NUL = uint8(iota)
    PROTO_TCP
    PROTO_UDP
)

var StrToProto = map[string]uint8{
    "tcp": PROTO_TCP,
    "udp": PROTO_UDP,
}

var ProtoToStr = map[uint8]string{
    PROTO_TCP: "tcp",
    PROTO_UDP: "udp",
}

// Actions
const (
    ACTION_SEND_DATA   = uint8(iota)
    ACTION_CLOSE_READ
    ACTION_CLOSE_WRITE
    ACTION_CLOSE_ALL = ACTION_CLOSE_READ | ACTION_CLOSE_WRITE
)