// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/type/date.proto

package date

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

// Represents a whole calendar date, e.g. date of birth. The time of day and
// time zone are either specified elsewhere or are not significant. The date
// is relative to the Proleptic Gregorian Calendar. The day may be 0 to
// represent a year and month where the day is not significant, e.g. credit card
// expiration date. The year may be 0 to represent a month and day independent
// of year, e.g. anniversary date. Related types are [google.type.TimeOfDay][google.type.TimeOfDay]
// and `google.protobuf.Timestamp`.
type Date struct {
	// Year of date. Must be from 1 to 9999, or 0 if specifying a date without
	// a year.
	Year int32 `protobuf:"varint,1,opt,name=year,proto3" json:"year,omitempty"`
	// Month of year. Must be from 1 to 12.
	Month int32 `protobuf:"varint,2,opt,name=month,proto3" json:"month,omitempty"`
	// Day of month. Must be from 1 to 31 and valid for the year and month, or 0
	// if specifying a year/month where the day is not significant.
	Day                  int32    `protobuf:"varint,3,opt,name=day,proto3" json:"day,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Date) Reset()         { *m = Date{} }
func (m *Date) String() string { return proto.CompactTextString(m) }
func (*Date) ProtoMessage()    {}
func (*Date) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c30699df886e3f, []int{0}
}

func (m *Date) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Date.Unmarshal(m, b)
}
func (m *Date) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Date.Marshal(b, m, deterministic)
}
func (m *Date) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Date.Merge(m, src)
}
func (m *Date) XXX_Size() int {
	return xxx_messageInfo_Date.Size(m)
}
func (m *Date) XXX_DiscardUnknown() {
	xxx_messageInfo_Date.DiscardUnknown(m)
}

var xxx_messageInfo_Date proto.InternalMessageInfo

func (m *Date) GetYear() int32 {
	if m != nil {
		return m.Year
	}
	return 0
}

func (m *Date) GetMonth() int32 {
	if m != nil {
		return m.Month
	}
	return 0
}

func (m *Date) GetDay() int32 {
	if m != nil {
		return m.Day
	}
	return 0
}

func init() {
	proto.RegisterType((*Date)(nil), "google.type.Date")
}

func init() { proto.RegisterFile("google/type/date.proto", fileDescriptor_92c30699df886e3f) }

var fileDescriptor_92c30699df886e3f = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4b, 0xcf, 0xcf, 0x4f,
	0xcf, 0x49, 0xd5, 0x2f, 0xa9, 0x2c, 0x48, 0xd5, 0x4f, 0x49, 0x2c, 0x49, 0xd5, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0xe2, 0x86, 0x88, 0xeb, 0x81, 0xc4, 0x95, 0x9c, 0xb8, 0x58, 0x5c, 0x12, 0x4b,
	0x52, 0x85, 0x84, 0xb8, 0x58, 0x2a, 0x53, 0x13, 0x8b, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83,
	0xc0, 0x6c, 0x21, 0x11, 0x2e, 0xd6, 0xdc, 0xfc, 0xbc, 0x92, 0x0c, 0x09, 0x26, 0xb0, 0x20, 0x84,
	0x23, 0x24, 0xc0, 0xc5, 0x9c, 0x92, 0x58, 0x29, 0xc1, 0x0c, 0x16, 0x03, 0x31, 0x9d, 0x62, 0xb9,
	0xf8, 0x93, 0xf3, 0x73, 0xf5, 0x90, 0x8c, 0x75, 0xe2, 0x04, 0x19, 0x1a, 0x00, 0xb2, 0x2e, 0x80,
	0x31, 0xca, 0x04, 0x2a, 0x93, 0x9e, 0x9f, 0x93, 0x98, 0x97, 0xae, 0x97, 0x5f, 0x94, 0xae, 0x9f,
	0x9e, 0x9a, 0x07, 0x76, 0x8c, 0x3e, 0x44, 0x2a, 0xb1, 0x20, 0xb3, 0x18, 0xe1, 0x4e, 0x6b, 0x10,
	0xf1, 0x83, 0x91, 0x71, 0x11, 0x13, 0xb3, 0x7b, 0x48, 0x40, 0x12, 0x1b, 0x58, 0xa5, 0x31, 0x20,
	0x00, 0x00, 0xff, 0xff, 0x84, 0x95, 0xf3, 0x4c, 0xd0, 0x00, 0x00, 0x00,
}
