// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/pubkey.proto

package echalotte_pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The type of encryption key.
type KeyType int32

const (
	KeyType_Curve25519 KeyType = 0
)

var KeyType_name = map[int32]string{
	0: "Curve25519",
}
var KeyType_value = map[string]int32{
	"Curve25519": 0,
}

func (x KeyType) String() string {
	return proto.EnumName(KeyType_name, int32(x))
}
func (KeyType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_pubkey_d4b2e54b605754f5, []int{0}
}

// An encryption public key.
type PublicKey struct {
	Type                 KeyType              `protobuf:"varint,1,opt,name=type,proto3,enum=echalotte.pb.KeyType" json:"type,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Data                 []byte               `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	SignatureKey         []byte               `protobuf:"bytes,10,opt,name=signature_key,json=signatureKey,proto3" json:"signature_key,omitempty"`
	Signature            []byte               `protobuf:"bytes,11,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *PublicKey) Reset()         { *m = PublicKey{} }
func (m *PublicKey) String() string { return proto.CompactTextString(m) }
func (*PublicKey) ProtoMessage()    {}
func (*PublicKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_pubkey_d4b2e54b605754f5, []int{0}
}
func (m *PublicKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublicKey.Unmarshal(m, b)
}
func (m *PublicKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublicKey.Marshal(b, m, deterministic)
}
func (dst *PublicKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublicKey.Merge(dst, src)
}
func (m *PublicKey) XXX_Size() int {
	return xxx_messageInfo_PublicKey.Size(m)
}
func (m *PublicKey) XXX_DiscardUnknown() {
	xxx_messageInfo_PublicKey.DiscardUnknown(m)
}

var xxx_messageInfo_PublicKey proto.InternalMessageInfo

func (m *PublicKey) GetType() KeyType {
	if m != nil {
		return m.Type
	}
	return KeyType_Curve25519
}

func (m *PublicKey) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *PublicKey) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *PublicKey) GetSignatureKey() []byte {
	if m != nil {
		return m.SignatureKey
	}
	return nil
}

func (m *PublicKey) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*PublicKey)(nil), "echalotte.pb.PublicKey")
	proto.RegisterEnum("echalotte.pb.KeyType", KeyType_name, KeyType_value)
}

func init() { proto.RegisterFile("pb/pubkey.proto", fileDescriptor_pubkey_d4b2e54b605754f5) }

var fileDescriptor_pubkey_d4b2e54b605754f5 = []byte{
	// 238 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8e, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0x8d, 0x0e, 0xa5, 0x6f, 0x75, 0x4a, 0x40, 0x88, 0x43, 0xb0, 0xe8, 0xa5, 0x7a, 0x48,
	0xb1, 0xb2, 0xc3, 0x8e, 0xe2, 0x71, 0x17, 0x29, 0xbb, 0x8f, 0xa4, 0x7b, 0xd6, 0xb2, 0xce, 0x84,
	0xec, 0x45, 0xc8, 0x7f, 0xe8, 0x9f, 0x25, 0x66, 0x6b, 0xf5, 0x96, 0x7c, 0xdf, 0x8f, 0xf7, 0xfd,
	0xe0, 0xc2, 0xea, 0xc2, 0x7a, 0xbd, 0xc1, 0x20, 0xad, 0x33, 0x64, 0x78, 0x8a, 0xf5, 0x87, 0xea,
	0x0c, 0x11, 0x4a, 0xab, 0xa7, 0xb7, 0x8d, 0x31, 0x4d, 0x87, 0x45, 0xec, 0xb4, 0x7f, 0x2f, 0xa8,
	0xdd, 0xe2, 0x8e, 0xd4, 0xd6, 0xee, 0xf1, 0xbb, 0x6f, 0x06, 0xc9, 0x9b, 0xd7, 0x5d, 0x5b, 0x2f,
	0x30, 0xf0, 0x07, 0x18, 0x51, 0xb0, 0x28, 0x58, 0xc6, 0xf2, 0x49, 0x79, 0x25, 0xff, 0xdf, 0x92,
	0x0b, 0x0c, 0xcb, 0x60, 0xb1, 0x8a, 0x08, 0x9f, 0x03, 0xd4, 0x0e, 0x15, 0xe1, 0x7a, 0xa5, 0x48,
	0x1c, 0x67, 0x2c, 0x1f, 0x97, 0x53, 0xb9, 0x9f, 0x93, 0xfd, 0x9c, 0x5c, 0xf6, 0x73, 0x55, 0x72,
	0xa0, 0x5f, 0x88, 0x73, 0x18, 0xad, 0x15, 0x29, 0x71, 0x92, 0xb1, 0x3c, 0xad, 0xe2, 0x9b, 0xdf,
	0xc3, 0xf9, 0xae, 0x6d, 0x3e, 0x15, 0x79, 0x87, 0xab, 0x0d, 0x06, 0x01, 0xb1, 0x4c, 0x87, 0xf0,
	0x57, 0xef, 0x06, 0x92, 0xe1, 0x2f, 0xc6, 0x11, 0xf8, 0x0b, 0x1e, 0xaf, 0xe1, 0xec, 0xa0, 0xc8,
	0x27, 0x00, 0xaf, 0xde, 0x7d, 0x61, 0x39, 0x9b, 0x3d, 0xcd, 0x2f, 0x8f, 0xf4, 0x69, 0x14, 0x7a,
	0xfe, 0x09, 0x00, 0x00, 0xff, 0xff, 0xa9, 0xb5, 0x23, 0xbd, 0x2e, 0x01, 0x00, 0x00,
}