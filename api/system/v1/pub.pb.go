// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/system/v1/pub.proto

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

func init() { proto.RegisterFile("api/system/v1/pub.proto", fileDescriptor_6145852d37fe055b) }

var fileDescriptor_6145852d37fe055b = []byte{
	// 363 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0x86, 0xbf, 0x7c, 0x3f, 0xfd, 0xea, 0x54, 0x45, 0x07, 0x41, 0x1a, 0x4b, 0x90, 0x2c, 0x0b,
	0x66, 0xa8, 0xee, 0xea, 0xce, 0xa2, 0x6e, 0xba, 0x28, 0x41, 0x05, 0x5d, 0x14, 0xa6, 0xed, 0x98,
	0x0e, 0xa6, 0x99, 0x31, 0x73, 0x12, 0xed, 0xd6, 0x5b, 0x70, 0xe3, 0xda, 0xab, 0x71, 0x29, 0x78,
	0x03, 0x52, 0xbd, 0x10, 0x99, 0x49, 0x0b, 0x26, 0xd6, 0x9f, 0xed, 0xfb, 0x9c, 0xf3, 0x30, 0x9c,
	0x79, 0xd1, 0x3a, 0x95, 0x9c, 0xa8, 0xb1, 0x02, 0x36, 0x22, 0x69, 0x83, 0xc8, 0xa4, 0xe7, 0xc9,
	0x58, 0x80, 0xc0, 0x4b, 0x54, 0x72, 0x2f, 0x03, 0x5e, 0xda, 0xb0, 0xed, 0xfc, 0xdc, 0x14, 0x98,
	0x51, 0xbb, 0x16, 0x08, 0x11, 0x84, 0x8c, 0xe8, 0x11, 0x1a, 0x45, 0x02, 0x28, 0x70, 0x11, 0xa9,
	0x8c, 0x6e, 0xdf, 0xff, 0x45, 0x7f, 0x3a, 0x49, 0x0f, 0x9f, 0xa0, 0xf2, 0x21, 0x8b, 0x58, 0x4c,
	0x81, 0xe1, 0x0d, 0x2f, 0x67, 0xf7, 0xf6, 0x47, 0x12, 0xc6, 0x3e, 0xbb, 0x4c, 0x98, 0x02, 0xbb,
	0x08, 0x5b, 0x54, 0x42, 0x7f, 0x48, 0x7d, 0x26, 0xc3, 0xb1, 0xbb, 0x72, 0xf3, 0xf4, 0x7a, 0xfb,
	0x1b, 0xe1, 0x32, 0xe9, 0x67, 0x31, 0xee, 0xa2, 0xb2, 0xcf, 0x02, 0xae, 0x80, 0xc5, 0xd8, 0x29,
	0xac, 0xce, 0xc0, 0x4c, 0x5d, 0xfb, 0x94, 0x6b, 0xf7, 0x9a, 0x71, 0x2f, 0xbb, 0x0b, 0x24, 0x9e,
	0xe6, 0x4d, 0xab, 0x8e, 0x7d, 0xf4, 0xaf, 0x2d, 0x02, 0x1e, 0x7d, 0x78, 0xb4, 0x49, 0x67, 0xe6,
	0xea, 0x7c, 0xa8, 0xb5, 0xab, 0x46, 0x5b, 0x71, 0x4b, 0x24, 0xd4, 0xa1, 0x76, 0x1e, 0xa1, 0x52,
	0x5b, 0x04, 0x22, 0x81, 0xaf, 0x2f, 0x51, 0x9d, 0x0f, 0xb5, 0x14, 0x1b, 0xe9, 0xa2, 0xfb, 0x5f,
	0x4b, 0x45, 0x02, 0xda, 0xda, 0x45, 0x15, 0x9f, 0x41, 0xcc, 0x59, 0xca, 0x3a, 0x57, 0x83, 0x39,
	0xc7, 0xc8, 0xd8, 0x0f, 0xec, 0xef, 0x2f, 0x91, 0x2d, 0x69, 0xff, 0x29, 0x42, 0xc7, 0x32, 0x14,
	0x74, 0x70, 0xc0, 0xc3, 0x6f, 0xfe, 0xd0, 0x2e, 0xc0, 0x6c, 0xaf, 0xf8, 0xf4, 0xc4, 0xa4, 0x4d,
	0xab, 0xbe, 0xd7, 0x7a, 0x98, 0x38, 0xd6, 0xe3, 0xc4, 0xb1, 0x9e, 0x27, 0x8e, 0x75, 0xf7, 0xe2,
	0xfc, 0x42, 0xf9, 0xfe, 0x75, 0xac, 0xb3, 0xcd, 0x8b, 0x98, 0x82, 0x50, 0xd7, 0x5b, 0xe7, 0x54,
	0x0d, 0xb9, 0x88, 0x48, 0xae, 0x91, 0xbb, 0x69, 0xa3, 0x57, 0x32, 0x85, 0xdb, 0x79, 0x0b, 0x00,
	0x00, 0xff, 0xff, 0x69, 0xed, 0x3e, 0x09, 0xd4, 0x02, 0x00, 0x00,
}
