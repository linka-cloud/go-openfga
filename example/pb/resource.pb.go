// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: example/pb/resource.proto

package resource

import (
	_ "github.com/alta/protopatch/patch/gopb"
	_ "go.linka.cloud/go-openfga/openfga"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_example_pb_resource_proto protoreflect.FileDescriptor

var file_example_pb_resource_proto_rawDesc = []byte{
	0x0a, 0x19, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x15, 0x6f, 0x70, 0x65, 0x6e, 0x66, 0x67, 0x61, 0x2f, 0x6f,
	0x70, 0x65, 0x6e, 0x66, 0x67, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x70, 0x61,
	0x74, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x62, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x32, 0xbf, 0x0b, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x12, 0x17, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x92, 0xfa, 0x06, 0x13, 0x0a, 0x11, 0x22, 0x0f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x12, 0x68,
	0x0a, 0x04, 0x52, 0x65, 0x61, 0x64, 0x12, 0x15, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x31, 0x92, 0xfa, 0x06, 0x2d, 0x0a, 0x11, 0x22, 0x0f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x0a, 0x18,
	0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x04, 0x7b, 0x69, 0x64, 0x7d,
	0x22, 0x06, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x76, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x17, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x39, 0x92, 0xfa, 0x06, 0x35, 0x0a, 0x11, 0x22, 0x0f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x0a, 0x20,
	0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x0d, 0x7b, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x69, 0x64, 0x7d, 0x22, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x12, 0x5a, 0x0a, 0x06, 0x41, 0x64, 0x64, 0x53, 0x75, 0x62, 0x12, 0x17, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x41, 0x64, 0x64, 0x53, 0x75, 0x62, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x41,
	0x64, 0x64, 0x53, 0x75, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x92,
	0xfa, 0x06, 0x19, 0x0a, 0x17, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x04, 0x7b, 0x69, 0x64, 0x7d, 0x22, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x7e, 0x0a, 0x07,
	0x52, 0x65, 0x61, 0x64, 0x53, 0x75, 0x62, 0x12, 0x18, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x53, 0x75, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x61,
	0x64, 0x53, 0x75, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3e, 0x92, 0xfa,
	0x06, 0x3a, 0x0a, 0x23, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x0d,
	0x7b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x7d, 0x28, 0x01, 0x22,
	0x06, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x0a, 0x13, 0x0a, 0x03, 0x73, 0x75, 0x62, 0x12, 0x04,
	0x7b, 0x69, 0x64, 0x7d, 0x22, 0x06, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x6d, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x17, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x18, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x30, 0x92, 0xfa, 0x06, 0x2c, 0x0a,
	0x11, 0x22, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x72, 0x0a, 0x17, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x04,
	0x7b, 0x69, 0x64, 0x7d, 0x22, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x4e, 0x0a, 0x04, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x15, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x17, 0x92, 0xfa, 0x06, 0x13, 0x0a, 0x11, 0x22, 0x0f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x4c, 0x0a, 0x05, 0x57,
	0x61, 0x74, 0x63, 0x68, 0x12, 0x16, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x18, 0x92,
	0xfa, 0x06, 0x14, 0x0a, 0x12, 0x22, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f,
	0x77, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x30, 0x01, 0x1a, 0x8a, 0x05, 0x92, 0xfa, 0x06, 0xf0,
	0x04, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0xe5, 0x01, 0x0a, 0x06,
	0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x3e, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x2c, 0x5b, 0x75, 0x73, 0x65, 0x72, 0x2c,
	0x20, 0x75, 0x73, 0x65, 0x72, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x6e, 0x6f, 0x6e, 0x5f, 0x65,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x5f, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x5d, 0x20, 0x6f, 0x72,
	0x20, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x2b, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x12, 0x18, 0x5b, 0x75, 0x73, 0x65, 0x72,
	0x5d, 0x20, 0x6f, 0x72, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x12, 0x35, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f,
	0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x22, 0x5b, 0x75, 0x73, 0x65, 0x72, 0x5d, 0x20, 0x6f,
	0x72, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x20, 0x6f, 0x72, 0x20, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x10, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x77, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x23,
	0x5b, 0x75, 0x73, 0x65, 0x72, 0x5d, 0x20, 0x6f, 0x72, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x20, 0x6f, 0x72, 0x20, 0x77, 0x61, 0x74, 0x63,
	0x68, 0x65, 0x72, 0x1a, 0x7e, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x12, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x08, 0x5b, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x5d, 0x12, 0x2d, 0x0a, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x24, 0x5b, 0x75,
	0x73, 0x65, 0x72, 0x5d, 0x20, 0x6f, 0x72, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x12, 0x2f, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x25, 0x5b, 0x75,
	0x73, 0x65, 0x72, 0x5d, 0x20, 0x6f, 0x72, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5f, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x73, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x1a, 0x6f, 0x0a, 0x03, 0x73, 0x75, 0x62, 0x12, 0x16, 0x0a, 0x08, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x0a, 0x5b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5d, 0x12, 0x26, 0x0a, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x1d, 0x5b, 0x75, 0x73,
	0x65, 0x72, 0x5d, 0x20, 0x6f, 0x72, 0x20, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x20, 0x66, 0x72, 0x6f,
	0x6d, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x72, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x12, 0x1e, 0x5b, 0x75, 0x73, 0x65, 0x72, 0x5d, 0x20, 0x6f, 0x72, 0x20,
	0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x22, 0x8a, 0x01, 0x6e, 0x6f, 0x6e, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x65, 0x64, 0x5f, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x28, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x3a, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2c, 0x20, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x3a, 0x20, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2c, 0x20, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x20, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x29, 0x20, 0x7b, 0x20, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x20, 0x3c, 0x20, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x20, 0x2b,
	0x20, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20,
	0x7d, 0x9a, 0xfa, 0x06, 0x11, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x07, 0x64,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x42, 0x6f, 0xca, 0xb5, 0x03, 0x02, 0x08, 0x01, 0x0a, 0x0c,
	0x63, 0x6f, 0x6d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x0d, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x0a, 0x2e,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0xa2, 0x02, 0x03, 0x52, 0x58, 0x58, 0xaa,
	0x02, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0xca, 0x02, 0x08, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0xe2, 0x02, 0x14, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x08, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_example_pb_resource_proto_goTypes = []interface{}{
	(*CreateRequest)(nil),   // 0: resource.CreateRequest
	(*ReadRequest)(nil),     // 1: resource.ReadRequest
	(*UpdateRequest)(nil),   // 2: resource.UpdateRequest
	(*AddSubRequest)(nil),   // 3: resource.AddSubRequest
	(*ReadSubRequest)(nil),  // 4: resource.ReadSubRequest
	(*DeleteRequest)(nil),   // 5: resource.DeleteRequest
	(*ListRequest)(nil),     // 6: resource.ListRequest
	(*WatchRequest)(nil),    // 7: resource.WatchRequest
	(*CreateResponse)(nil),  // 8: resource.CreateResponse
	(*ReadResponse)(nil),    // 9: resource.ReadResponse
	(*UpdateResponse)(nil),  // 10: resource.UpdateResponse
	(*AddSubResponse)(nil),  // 11: resource.AddSubResponse
	(*ReadSubResponse)(nil), // 12: resource.ReadSubResponse
	(*DeleteResponse)(nil),  // 13: resource.DeleteResponse
	(*ListResponse)(nil),    // 14: resource.ListResponse
	(*Event)(nil),           // 15: resource.Event
}
var file_example_pb_resource_proto_depIdxs = []int32{
	0,  // 0: resource.ResourceService.Create:input_type -> resource.CreateRequest
	1,  // 1: resource.ResourceService.Read:input_type -> resource.ReadRequest
	2,  // 2: resource.ResourceService.Update:input_type -> resource.UpdateRequest
	3,  // 3: resource.ResourceService.AddSub:input_type -> resource.AddSubRequest
	4,  // 4: resource.ResourceService.ReadSub:input_type -> resource.ReadSubRequest
	5,  // 5: resource.ResourceService.Delete:input_type -> resource.DeleteRequest
	6,  // 6: resource.ResourceService.List:input_type -> resource.ListRequest
	7,  // 7: resource.ResourceService.Watch:input_type -> resource.WatchRequest
	8,  // 8: resource.ResourceService.Create:output_type -> resource.CreateResponse
	9,  // 9: resource.ResourceService.Read:output_type -> resource.ReadResponse
	10, // 10: resource.ResourceService.Update:output_type -> resource.UpdateResponse
	11, // 11: resource.ResourceService.AddSub:output_type -> resource.AddSubResponse
	12, // 12: resource.ResourceService.ReadSub:output_type -> resource.ReadSubResponse
	13, // 13: resource.ResourceService.Delete:output_type -> resource.DeleteResponse
	14, // 14: resource.ResourceService.List:output_type -> resource.ListResponse
	15, // 15: resource.ResourceService.Watch:output_type -> resource.Event
	8,  // [8:16] is the sub-list for method output_type
	0,  // [0:8] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_example_pb_resource_proto_init() }
func file_example_pb_resource_proto_init() {
	if File_example_pb_resource_proto != nil {
		return
	}
	file_example_pb_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_example_pb_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_example_pb_resource_proto_goTypes,
		DependencyIndexes: file_example_pb_resource_proto_depIdxs,
	}.Build()
	File_example_pb_resource_proto = out.File
	file_example_pb_resource_proto_rawDesc = nil
	file_example_pb_resource_proto_goTypes = nil
	file_example_pb_resource_proto_depIdxs = nil
}
