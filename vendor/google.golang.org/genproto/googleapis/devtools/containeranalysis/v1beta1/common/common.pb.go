// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/devtools/containeranalysis/v1beta1/common/common.proto

package common

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
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

// Kind represents the kinds of notes supported.
type NoteKind int32

const (
	// Unknown.
	NoteKind_NOTE_KIND_UNSPECIFIED NoteKind = 0
	// The note and occurrence represent a package vulnerability.
	NoteKind_VULNERABILITY NoteKind = 1
	// The note and occurrence assert build provenance.
	NoteKind_BUILD NoteKind = 2
	// This represents an image basis relationship.
	NoteKind_IMAGE NoteKind = 3
	// This represents a package installed via a package manager.
	NoteKind_PACKAGE NoteKind = 4
	// The note and occurrence track deployment events.
	NoteKind_DEPLOYMENT NoteKind = 5
	// The note and occurrence track the initial discovery status of a resource.
	NoteKind_DISCOVERY NoteKind = 6
	// This represents a logical "role" that can attest to artifacts.
	NoteKind_ATTESTATION NoteKind = 7
)

var NoteKind_name = map[int32]string{
	0: "NOTE_KIND_UNSPECIFIED",
	1: "VULNERABILITY",
	2: "BUILD",
	3: "IMAGE",
	4: "PACKAGE",
	5: "DEPLOYMENT",
	6: "DISCOVERY",
	7: "ATTESTATION",
}

var NoteKind_value = map[string]int32{
	"NOTE_KIND_UNSPECIFIED": 0,
	"VULNERABILITY":         1,
	"BUILD":                 2,
	"IMAGE":                 3,
	"PACKAGE":               4,
	"DEPLOYMENT":            5,
	"DISCOVERY":             6,
	"ATTESTATION":           7,
}

func (x NoteKind) String() string {
	return proto.EnumName(NoteKind_name, int32(x))
}

func (NoteKind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_baf5189c00300310, []int{0}
}

// Metadata for any related URL information.
type RelatedUrl struct {
	// Specific URL associated with the resource.
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// Label to describe usage of the URL.
	Label                string   `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RelatedUrl) Reset()         { *m = RelatedUrl{} }
func (m *RelatedUrl) String() string { return proto.CompactTextString(m) }
func (*RelatedUrl) ProtoMessage()    {}
func (*RelatedUrl) Descriptor() ([]byte, []int) {
	return fileDescriptor_baf5189c00300310, []int{0}
}

func (m *RelatedUrl) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelatedUrl.Unmarshal(m, b)
}
func (m *RelatedUrl) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelatedUrl.Marshal(b, m, deterministic)
}
func (m *RelatedUrl) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelatedUrl.Merge(m, src)
}
func (m *RelatedUrl) XXX_Size() int {
	return xxx_messageInfo_RelatedUrl.Size(m)
}
func (m *RelatedUrl) XXX_DiscardUnknown() {
	xxx_messageInfo_RelatedUrl.DiscardUnknown(m)
}

var xxx_messageInfo_RelatedUrl proto.InternalMessageInfo

func (m *RelatedUrl) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *RelatedUrl) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func init() {
	proto.RegisterEnum("grafeas.v1beta1.NoteKind", NoteKind_name, NoteKind_value)
	proto.RegisterType((*RelatedUrl)(nil), "grafeas.v1beta1.RelatedUrl")
}

func init() {
	proto.RegisterFile("google/devtools/containeranalysis/v1beta1/common/common.proto", fileDescriptor_baf5189c00300310)
}

var fileDescriptor_baf5189c00300310 = []byte{
	// 322 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0x41, 0x6b, 0xea, 0x40,
	0x14, 0x85, 0x5f, 0xf4, 0xa9, 0xcf, 0x2b, 0x3e, 0xa7, 0x43, 0x0b, 0xba, 0x2b, 0x5d, 0x95, 0x2e,
	0x12, 0xa4, 0xdd, 0x95, 0x2e, 0xa2, 0x99, 0xca, 0x60, 0x9c, 0x84, 0x38, 0x11, 0xec, 0x46, 0x46,
	0x9d, 0x0e, 0x81, 0x71, 0x46, 0x92, 0x54, 0x28, 0xfd, 0x09, 0xfd, 0x17, 0xfd, 0xa5, 0x45, 0xe3,
	0xaa, 0xab, 0xae, 0xce, 0x39, 0xf7, 0x5e, 0x2e, 0x9c, 0x0f, 0x9e, 0x94, 0xb5, 0x4a, 0x4b, 0x6f,
	0x2b, 0x0f, 0xa5, 0xb5, 0xba, 0xf0, 0x36, 0xd6, 0x94, 0x22, 0x33, 0x32, 0x17, 0x46, 0xe8, 0xf7,
	0x22, 0x2b, 0xbc, 0xc3, 0x70, 0x2d, 0x4b, 0x31, 0xf4, 0x36, 0x76, 0xb7, 0xb3, 0xe6, 0x2c, 0xee,
	0x3e, 0xb7, 0xa5, 0xc5, 0x3d, 0x95, 0x8b, 0x57, 0x29, 0x0a, 0xf7, 0x7c, 0x74, 0xf3, 0x00, 0x90,
	0x48, 0x2d, 0x4a, 0xb9, 0x4d, 0x73, 0x8d, 0x11, 0xd4, 0xdf, 0x72, 0xdd, 0x77, 0xae, 0x9d, 0xdb,
	0x76, 0x72, 0xb4, 0xf8, 0x12, 0x1a, 0x5a, 0xac, 0xa5, 0xee, 0xd7, 0x4e, 0xb3, 0x2a, 0xdc, 0x7d,
	0x3a, 0xf0, 0x8f, 0xd9, 0x52, 0x4e, 0x33, 0xb3, 0xc5, 0x03, 0xb8, 0x62, 0x11, 0x27, 0xab, 0x29,
	0x65, 0xc1, 0x2a, 0x65, 0xf3, 0x98, 0x8c, 0xe9, 0x33, 0x25, 0x01, 0xfa, 0x83, 0x2f, 0xa0, 0xbb,
	0x48, 0x43, 0x46, 0x12, 0x7f, 0x44, 0x43, 0xca, 0x97, 0xc8, 0xc1, 0x6d, 0x68, 0x8c, 0x52, 0x1a,
	0x06, 0xa8, 0x76, 0xb4, 0x74, 0xe6, 0x4f, 0x08, 0xaa, 0xe3, 0x0e, 0xb4, 0x62, 0x7f, 0x3c, 0x3d,
	0x86, 0xbf, 0xf8, 0x3f, 0x40, 0x40, 0xe2, 0x30, 0x5a, 0xce, 0x08, 0xe3, 0xa8, 0x81, 0xbb, 0xd0,
	0x0e, 0xe8, 0x7c, 0x1c, 0x2d, 0x48, 0xb2, 0x44, 0x4d, 0xdc, 0x83, 0x8e, 0xcf, 0x39, 0x99, 0x73,
	0x9f, 0xd3, 0x88, 0xa1, 0xd6, 0xe8, 0x03, 0x06, 0x99, 0x75, 0x7f, 0x34, 0x73, 0xab, 0xde, 0xb1,
	0xf3, 0xb2, 0xa8, 0x90, 0xb9, 0xca, 0x6a, 0x61, 0x94, 0x6b, 0x73, 0xe5, 0x29, 0x69, 0x4e, 0x3c,
	0xbc, 0x6a, 0x25, 0xf6, 0x59, 0xf1, 0x7b, 0xa2, 0x8f, 0x95, 0x7c, 0xd5, 0xea, 0x93, 0xc4, 0x5f,
	0x37, 0x4f, 0x8f, 0xee, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xe2, 0x0a, 0x57, 0x94, 0x99, 0x01,
	0x00, 0x00,
}
