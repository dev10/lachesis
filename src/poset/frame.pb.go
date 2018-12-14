// Code generated by protoc-gen-go. DO NOT EDIT.
// source: frame.proto

package poset

import (
	fmt "fmt"
	peers "github.com/Fantom-foundation/go-lachesis/src/peers"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Frame struct {
	Round                int64            `protobuf:"varint,1,opt,name=Round,proto3" json:"Round,omitempty"`
	Roots                map[string]*Root `protobuf:"bytes,2,rep,name=Roots,proto3" json:"Roots,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Events               []*EventMessage  `protobuf:"bytes,3,rep,name=Events,proto3" json:"Events,omitempty"`
	Peers                []*peers.Peer    `protobuf:"bytes,4,rep,name=Peers,proto3" json:"Peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Frame) Reset()         { *m = Frame{} }
func (m *Frame) String() string { return proto.CompactTextString(m) }
func (*Frame) ProtoMessage()    {}
func (*Frame) Descriptor() ([]byte, []int) {
	return fileDescriptor_5379e2b825e15002, []int{0}
}

func (m *Frame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Frame.Unmarshal(m, b)
}
func (m *Frame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Frame.Marshal(b, m, deterministic)
}
func (m *Frame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Frame.Merge(m, src)
}
func (m *Frame) XXX_Size() int {
	return xxx_messageInfo_Frame.Size(m)
}
func (m *Frame) XXX_DiscardUnknown() {
	xxx_messageInfo_Frame.DiscardUnknown(m)
}

var xxx_messageInfo_Frame proto.InternalMessageInfo

func (m *Frame) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *Frame) GetRoots() map[string]*Root {
	if m != nil {
		return m.Roots
	}
	return nil
}

func (m *Frame) GetEvents() []*EventMessage {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *Frame) GetPeers() []*peers.Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

func init() {
	proto.RegisterType((*Frame)(nil), "poset.Frame")
	proto.RegisterMapType((map[string]*Root)(nil), "poset.Frame.RootsEntry")
}

func init() { proto.RegisterFile("frame.proto", fileDescriptor_5379e2b825e15002) }

var fileDescriptor_5379e2b825e15002 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x4f, 0xcd, 0x4a, 0xc4, 0x30,
	0x10, 0x26, 0xad, 0x5d, 0x30, 0xbd, 0x48, 0x14, 0x0c, 0x3d, 0x55, 0x4f, 0x0b, 0xd2, 0x14, 0xf4,
	0x22, 0x82, 0xc7, 0xee, 0x4d, 0x90, 0xbc, 0x41, 0xb6, 0xce, 0x76, 0x8b, 0xdb, 0x4c, 0x49, 0xd2,
	0x85, 0x7d, 0x5e, 0x5f, 0x44, 0x32, 0xa9, 0x78, 0x09, 0x99, 0xef, 0x67, 0xbe, 0x6f, 0x78, 0x79,
	0x70, 0x66, 0x02, 0x35, 0x3b, 0x0c, 0x28, 0x8a, 0x19, 0x3d, 0x84, 0x8a, 0x3b, 0xc4, 0x90, 0xa0,
	0xaa, 0x84, 0x33, 0xd8, 0xbf, 0xe1, 0x7d, 0x18, 0xc3, 0x71, 0xd9, 0xab, 0x1e, 0xa7, 0x76, 0x67,
	0x6c, 0xc0, 0xa9, 0x39, 0xe0, 0x62, 0xbf, 0x4c, 0x18, 0xd1, 0xb6, 0x03, 0x36, 0x27, 0xd3, 0x1f,
	0xc1, 0x8f, 0xbe, 0xf5, 0xae, 0x6f, 0x67, 0x00, 0xe7, 0xe9, 0x4d, 0xf6, 0xc7, 0x1f, 0xc6, 0x8b,
	0x5d, 0x8c, 0x13, 0x77, 0xbc, 0xd0, 0xd1, 0x28, 0x59, 0xcd, 0xb6, 0xb9, 0x4e, 0x83, 0x68, 0x22,
	0x8a, 0xc1, 0xcb, 0xac, 0xce, 0xb7, 0xe5, 0xf3, 0xbd, 0xa2, 0x3a, 0x8a, 0x2c, 0x8a, 0x98, 0xce,
	0x06, 0x77, 0xd1, 0x49, 0x25, 0x9e, 0xf8, 0xa6, 0x8b, 0xe5, 0xbc, 0xcc, 0x49, 0x7f, 0xbb, 0xea,
	0x09, 0xfc, 0x00, 0xef, 0xcd, 0x00, 0x7a, 0x95, 0x88, 0x07, 0x5e, 0x7c, 0xc6, 0x3e, 0xf2, 0x8a,
	0xb4, 0xa5, 0xa2, 0x76, 0x2a, 0x62, 0x3a, 0x31, 0x55, 0xc7, 0xf9, 0x7f, 0x88, 0xb8, 0xe1, 0xf9,
	0x37, 0x5c, 0xa8, 0xe0, 0xb5, 0x8e, 0xdf, 0xb8, 0xe2, 0x6c, 0x4e, 0x0b, 0xc8, 0xac, 0x66, 0x69,
	0x05, 0xc5, 0x45, 0x8f, 0x4e, 0xcc, 0x5b, 0xf6, 0xca, 0xf6, 0x1b, 0x3a, 0xf6, 0xe5, 0x37, 0x00,
	0x00, 0xff, 0xff, 0xb0, 0x66, 0x3b, 0x84, 0x5a, 0x01, 0x00, 0x00,
}
