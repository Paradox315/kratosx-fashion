// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/system/v1/resource.proto

package v1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("api/system/v1/resource.proto", fileDescriptor_0bc31b3b1de9fdc1) }

var fileDescriptor_0bc31b3b1de9fdc1 = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0xcf, 0x4a, 0xfb, 0x40,
	0x10, 0xc7, 0x7f, 0xf9, 0xe1, 0xdf, 0xb5, 0x5a, 0xbb, 0x94, 0x5a, 0x63, 0x09, 0xd2, 0x63, 0xc1,
	0x2c, 0xd5, 0x9b, 0xde, 0xda, 0x4a, 0xf1, 0x20, 0x48, 0x50, 0x28, 0xbd, 0xad, 0x76, 0xac, 0xc1,
	0x34, 0x1b, 0x77, 0x27, 0xc5, 0x20, 0xbd, 0xf8, 0x0a, 0x5e, 0x7c, 0x05, 0xdf, 0xc4, 0xa3, 0xe0,
	0x0b, 0x48, 0xf5, 0x41, 0x24, 0x9b, 0x88, 0x4d, 0xed, 0x49, 0xbc, 0xce, 0x77, 0xe6, 0xf3, 0xc9,
	0x64, 0x87, 0x54, 0x78, 0xe0, 0x32, 0x15, 0x29, 0x84, 0x01, 0x1b, 0xd6, 0x99, 0x04, 0x25, 0x42,
	0x79, 0x01, 0x76, 0x20, 0x05, 0x0a, 0xba, 0xca, 0x03, 0xd7, 0x4e, 0x52, 0x7b, 0x58, 0x37, 0xcd,
	0x6c, 0x73, 0x1a, 0xe8, 0x56, 0xb3, 0xd2, 0x17, 0xa2, 0xef, 0x01, 0x8b, 0x5b, 0xb8, 0xef, 0x0b,
	0xe4, 0xe8, 0x0a, 0x5f, 0x25, 0xe9, 0xee, 0xd3, 0x1c, 0x59, 0x72, 0x52, 0x36, 0x75, 0x08, 0x69,
	0x4a, 0xe0, 0x08, 0xc7, 0xe0, 0x87, 0xd4, 0xb4, 0x33, 0x12, 0x3b, 0x2e, 0x3a, 0x70, 0x13, 0x82,
	0x42, 0xb3, 0x34, 0x95, 0x1d, 0xb5, 0x1c, 0x08, 0xbc, 0xa8, 0xba, 0x7e, 0xff, 0xfa, 0xf1, 0xf0,
	0x9f, 0x54, 0xe7, 0xd9, 0x00, 0xfc, 0x70, 0xdf, 0xa8, 0xc5, 0xcc, 0xb3, 0xa0, 0xf7, 0x07, 0x4c,
	0xf3, 0x9b, 0xd9, 0x21, 0xa4, 0x05, 0x1e, 0xa4, 0xcc, 0xcd, 0x1f, 0x73, 0xea, 0x0b, 0x39, 0x1d,
	0x1d, 0x0e, 0x02, 0x8c, 0x12, 0x6a, 0x51, 0x53, 0xd7, 0x6a, 0x39, 0x4d, 0x55, 0xec, 0xce, 0xed,
	0xa9, 0x11, 0xed, 0x92, 0x95, 0x36, 0x60, 0x8c, 0x3d, 0x95, 0x00, 0xb4, 0x3c, 0xe3, 0x93, 0x12,
	0x72, 0x79, 0xe6, 0x22, 0x31, 0x78, 0x43, 0x83, 0x0b, 0x34, 0xaf, 0xc1, 0x0c, 0x25, 0x40, 0x0c,
	0x1f, 0x51, 0x20, 0x85, 0x09, 0x76, 0x23, 0x72, 0x84, 0xf7, 0x3b, 0x43, 0x45, 0x1b, 0x4a, 0xb4,
	0x38, 0x61, 0x90, 0xc2, 0x4b, 0x35, 0x1d, 0x92, 0x6b, 0x03, 0x3a, 0x22, 0x44, 0xd0, 0x3b, 0x6c,
	0xcd, 0xfe, 0x07, 0x89, 0x64, 0xfa, 0x3d, 0xf4, 0x98, 0x4c, 0x34, 0x79, 0xad, 0x59, 0xa6, 0x8b,
	0x4c, 0xea, 0x6a, 0xa3, 0xf9, 0x3c, 0xb6, 0x8c, 0x97, 0xb1, 0x65, 0xbc, 0x8d, 0x2d, 0xe3, 0xf1,
	0xdd, 0xfa, 0x47, 0xb2, 0x67, 0x78, 0x62, 0x74, 0xb7, 0xaf, 0x25, 0x47, 0xa1, 0x6e, 0x77, 0x2e,
	0xb9, 0xba, 0x72, 0x85, 0xcf, 0x32, 0x87, 0x79, 0x30, 0xac, 0x9f, 0x2f, 0xe8, 0xbb, 0xdb, 0xfb,
	0x0c, 0x00, 0x00, 0xff, 0xff, 0x08, 0xb6, 0x52, 0xda, 0xe0, 0x02, 0x00, 0x00,
}
