// Code generated by protoc-gen-go.
// source: server.proto
// DO NOT EDIT!

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ServerInfo struct {
	Id               *int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name             *string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	CmdAddr          *string `protobuf:"bytes,4,opt,name=cmdAddr" json:"cmdAddr,omitempty"`
	RpcAddr          *string `protobuf:"bytes,5,opt,name=rpcAddr" json:"rpcAddr,omitempty"`
	ExtAddr          *string `protobuf:"bytes,6,opt,name=extAddr" json:"extAddr,omitempty"`
	LastTick         *int64  `protobuf:"varint,7,opt,name=lastTick" json:"lastTick,omitempty"`
	ServerType       *string `protobuf:"bytes,8,opt,name=serverType" json:"serverType,omitempty"`
	Load             *int32  `protobuf:"varint,9,opt,name=load" json:"load,omitempty"`
	Conf             *string `protobuf:"bytes,10,opt,name=conf" json:"conf,omitempty"`
	OnlineCount      *int32  `protobuf:"varint,11,opt,name=onlineCount" json:"onlineCount,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ServerInfo) Reset()                    { *m = ServerInfo{} }
func (m *ServerInfo) String() string            { return proto.CompactTextString(m) }
func (*ServerInfo) ProtoMessage()               {}
func (*ServerInfo) Descriptor() ([]byte, []int) { return fileDescriptor6, []int{0} }

func (m *ServerInfo) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *ServerInfo) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *ServerInfo) GetCmdAddr() string {
	if m != nil && m.CmdAddr != nil {
		return *m.CmdAddr
	}
	return ""
}

func (m *ServerInfo) GetRpcAddr() string {
	if m != nil && m.RpcAddr != nil {
		return *m.RpcAddr
	}
	return ""
}

func (m *ServerInfo) GetExtAddr() string {
	if m != nil && m.ExtAddr != nil {
		return *m.ExtAddr
	}
	return ""
}

func (m *ServerInfo) GetLastTick() int64 {
	if m != nil && m.LastTick != nil {
		return *m.LastTick
	}
	return 0
}

func (m *ServerInfo) GetServerType() string {
	if m != nil && m.ServerType != nil {
		return *m.ServerType
	}
	return ""
}

func (m *ServerInfo) GetLoad() int32 {
	if m != nil && m.Load != nil {
		return *m.Load
	}
	return 0
}

func (m *ServerInfo) GetConf() string {
	if m != nil && m.Conf != nil {
		return *m.Conf
	}
	return ""
}

func (m *ServerInfo) GetOnlineCount() int32 {
	if m != nil && m.OnlineCount != nil {
		return *m.OnlineCount
	}
	return 0
}

func init() {
	proto.RegisterType((*ServerInfo)(nil), "ServerInfo")
}

func init() { proto.RegisterFile("server.proto", fileDescriptor6) }

var fileDescriptor6 = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0xcd, 0x31, 0x0e, 0x82, 0x30,
	0x18, 0xc5, 0xf1, 0x54, 0x40, 0xe0, 0x41, 0xa2, 0xa9, 0xcb, 0x37, 0x12, 0x27, 0x26, 0xef, 0x60,
	0x9c, 0x9c, 0xe5, 0x02, 0xa4, 0x2d, 0x09, 0x11, 0xfa, 0x91, 0x52, 0x8d, 0xde, 0xca, 0x23, 0x1a,
	0x6b, 0xba, 0xfe, 0xde, 0x4b, 0xfe, 0xa8, 0x57, 0xe3, 0x9e, 0xc6, 0x9d, 0x16, 0xc7, 0x9e, 0x8f,
	0x1f, 0x01, 0xdc, 0x02, 0x5c, 0xed, 0xc0, 0x12, 0xd8, 0x8c, 0x9a, 0x44, 0x23, 0xda, 0x4c, 0xd6,
	0x48, 0x6d, 0x3f, 0x1b, 0x4a, 0x1a, 0xd1, 0x96, 0x72, 0x87, 0x5c, 0xcd, 0xfa, 0xac, 0xb5, 0xa3,
	0x34, 0x82, 0x5b, 0x54, 0x80, 0x2c, 0x82, 0x79, 0xf9, 0x00, 0xdb, 0x00, 0x7b, 0x14, 0x53, 0xbf,
	0xfa, 0x6e, 0x54, 0x77, 0xca, 0x1b, 0xd1, 0x26, 0x52, 0x02, 0xff, 0x7a, 0xf7, 0x5e, 0x0c, 0x15,
	0xe1, 0x55, 0x23, 0x9d, 0xb8, 0xd7, 0x54, 0xc6, 0xa8, 0x62, 0x3b, 0x10, 0xc2, 0x76, 0x40, 0xc5,
	0x76, 0x1a, 0xad, 0xb9, 0xf0, 0xc3, 0x7a, 0xaa, 0x7e, 0x97, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x18, 0x12, 0x05, 0xd0, 0xc1, 0x00, 0x00, 0x00,
}
