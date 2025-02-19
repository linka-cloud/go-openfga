// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: openfga/v1alpha1/openfga.proto

package openfgav1alpha1

import (
	_ "github.com/alta/protopatch/patch/gopb"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Module struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Extends    []*Type  `protobuf:"bytes,2,rep,name=extends,proto3" json:"extends,omitempty"`
	Types      []*Type  `protobuf:"bytes,3,rep,name=types,proto3" json:"types,omitempty"`
	Conditions []string `protobuf:"bytes,4,rep,name=conditions,proto3" json:"conditions,omitempty"`
}

func (x *Module) Reset() {
	*x = Module{}
	if protoimpl.UnsafeEnabled {
		mi := &file_openfga_v1alpha1_openfga_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Module) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Module) ProtoMessage() {}

func (x *Module) ProtoReflect() protoreflect.Message {
	mi := &file_openfga_v1alpha1_openfga_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Module.ProtoReflect.Descriptor instead.
func (*Module) Descriptor() ([]byte, []int) {
	return file_openfga_v1alpha1_openfga_proto_rawDescGZIP(), []int{0}
}

func (x *Module) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Module) GetExtends() []*Type {
	if x != nil {
		return x.Extends
	}
	return nil
}

func (x *Module) GetTypes() []*Type {
	if x != nil {
		return x.Types
	}
	return nil
}

func (x *Module) GetConditions() []string {
	if x != nil {
		return x.Conditions
	}
	return nil
}

type Type struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Relations []*Relation `protobuf:"bytes,2,rep,name=relations,proto3" json:"relations,omitempty"`
}

func (x *Type) Reset() {
	*x = Type{}
	if protoimpl.UnsafeEnabled {
		mi := &file_openfga_v1alpha1_openfga_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Type) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Type) ProtoMessage() {}

func (x *Type) ProtoReflect() protoreflect.Message {
	mi := &file_openfga_v1alpha1_openfga_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Type.ProtoReflect.Descriptor instead.
func (*Type) Descriptor() ([]byte, []int) {
	return file_openfga_v1alpha1_openfga_proto_rawDescGZIP(), []int{1}
}

func (x *Type) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Type) GetRelations() []*Relation {
	if x != nil {
		return x.Relations
	}
	return nil
}

type Relation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Relation string `protobuf:"bytes,2,opt,name=relation,proto3" json:"relation,omitempty"`
}

func (x *Relation) Reset() {
	*x = Relation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_openfga_v1alpha1_openfga_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Relation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Relation) ProtoMessage() {}

func (x *Relation) ProtoReflect() protoreflect.Message {
	mi := &file_openfga_v1alpha1_openfga_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Relation.ProtoReflect.Descriptor instead.
func (*Relation) Descriptor() ([]byte, []int) {
	return file_openfga_v1alpha1_openfga_proto_rawDescGZIP(), []int{2}
}

func (x *Relation) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Relation) GetRelation() string {
	if x != nil {
		return x.Relation
	}
	return ""
}

type Access struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Relation string `protobuf:"bytes,2,opt,name=relation,proto3" json:"relation,omitempty"`
	ID       string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Access) Reset() {
	*x = Access{}
	if protoimpl.UnsafeEnabled {
		mi := &file_openfga_v1alpha1_openfga_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Access) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Access) ProtoMessage() {}

func (x *Access) ProtoReflect() protoreflect.Message {
	mi := &file_openfga_v1alpha1_openfga_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Access.ProtoReflect.Descriptor instead.
func (*Access) Descriptor() ([]byte, []int) {
	return file_openfga_v1alpha1_openfga_proto_rawDescGZIP(), []int{3}
}

func (x *Access) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Access) GetRelation() string {
	if x != nil {
		return x.Relation
	}
	return ""
}

func (x *Access) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

var file_openfga_v1alpha1_openfga_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*Module)(nil),
		Field:         14242,
		Name:          "fga.v1alpha1.module",
		Tag:           "bytes,14242,opt,name=module",
		Filename:      "openfga/v1alpha1/openfga.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*Access)(nil),
		Field:         14242,
		Name:          "fga.v1alpha1.access",
		Tag:           "bytes,14242,opt,name=access",
		Filename:      "openfga/v1alpha1/openfga.proto",
	},
}

// Extension fields to descriptorpb.ServiceOptions.
var (
	// optional fga.v1alpha1.Module module = 14242;
	ExtModule = &file_openfga_v1alpha1_openfga_proto_extTypes[0]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional fga.v1alpha1.Access access = 14242;
	ExtAccess = &file_openfga_v1alpha1_openfga_proto_extTypes[1]
)

var File_openfga_v1alpha1_openfga_proto protoreflect.FileDescriptor

var file_openfga_v1alpha1_openfga_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x6f, 0x70, 0x65, 0x6e, 0x66, 0x67, 0x61, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x66, 0x67, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0c, 0x66, 0x67, 0x61, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x20,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x2f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xad, 0x01, 0x0a, 0x06, 0x4d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x17, 0xfa, 0x42, 0x14, 0x72, 0x12, 0x32, 0x10, 0x5e, 0x5b, 0x5e, 0x3a, 0x23,
	0x40, 0x5c, 0x73, 0x5d, 0x7b, 0x31, 0x2c, 0x35, 0x30, 0x7d, 0x24, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x2c, 0x0a, 0x07, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x66, 0x67, 0x61, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x73, 0x12,
	0x28, 0x0a, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x66, 0x67, 0x61, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6e,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x63,
	0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x69, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x2b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x17, 0xfa, 0x42, 0x14, 0x72, 0x12, 0x32, 0x10, 0x5e, 0x5b, 0x5e, 0x3a, 0x23, 0x40, 0x5c, 0x73,
	0x5d, 0x7b, 0x31, 0x2c, 0x35, 0x30, 0x7d, 0x24, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34,
	0x0a, 0x09, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x66, 0x67, 0x61, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x22, 0x53, 0x0a, 0x08, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x2b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17,
	0xfa, 0x42, 0x14, 0x72, 0x12, 0x32, 0x10, 0x5e, 0x5b, 0x5e, 0x3a, 0x23, 0x40, 0x5c, 0x73, 0x5d,
	0x7b, 0x31, 0x2c, 0x35, 0x30, 0x7d, 0x24, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x95, 0x01, 0x0a, 0x06, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x17, 0xfa, 0x42, 0x14, 0x72, 0x12, 0x32, 0x10, 0x5e, 0x5b, 0x5e, 0x3a, 0x23,
	0x40, 0x5c, 0x73, 0x5d, 0x7b, 0x31, 0x2c, 0x35, 0x30, 0x7d, 0x24, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x34, 0x0a, 0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x18, 0xfa, 0x42, 0x15, 0x72, 0x13, 0x32, 0x11, 0x5e, 0x5b, 0x5e, 0x3a,
	0x23, 0x40, 0x5c, 0x73, 0x5d, 0x7b, 0x31, 0x2c, 0x32, 0x35, 0x34, 0x7d, 0x24, 0x52, 0x08, 0x72,
	0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x18, 0xfa, 0x42, 0x15, 0x72, 0x13, 0x32, 0x11, 0x5e, 0x5b, 0x5e, 0x3a,
	0x23, 0x40, 0x5c, 0x73, 0x5d, 0x7b, 0x31, 0x2c, 0x32, 0x35, 0x34, 0x7d, 0x24, 0x52, 0x02, 0x69,
	0x64, 0x3a, 0x4e, 0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa2, 0x6f, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x66, 0x67, 0x61, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x3a, 0x4d, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1e, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa2, 0x6f, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x66, 0x67, 0x61, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x42, 0xb3, 0x01, 0xca, 0xb5, 0x03, 0x02, 0x08, 0x01, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x66,
	0x67, 0x61, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x42, 0x0c, 0x4f, 0x70, 0x65,
	0x6e, 0x66, 0x67, 0x61, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3a, 0x67, 0x6f, 0x2e,
	0x6c, 0x69, 0x6e, 0x6b, 0x61, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x67, 0x6f, 0x2d, 0x6f,
	0x70, 0x65, 0x6e, 0x66, 0x67, 0x61, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x66, 0x67, 0x61, 0x2f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x6f, 0x70, 0x65, 0x6e, 0x66, 0x67, 0x61, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa, 0x02, 0x0c,
	0x46, 0x67, 0x61, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xca, 0x02, 0x0c, 0x46,
	0x67, 0x61, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x18, 0x46, 0x67,
	0x61, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x46, 0x67, 0x61, 0x3a, 0x3a, 0x56, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_openfga_v1alpha1_openfga_proto_rawDescOnce sync.Once
	file_openfga_v1alpha1_openfga_proto_rawDescData = file_openfga_v1alpha1_openfga_proto_rawDesc
)

func file_openfga_v1alpha1_openfga_proto_rawDescGZIP() []byte {
	file_openfga_v1alpha1_openfga_proto_rawDescOnce.Do(func() {
		file_openfga_v1alpha1_openfga_proto_rawDescData = protoimpl.X.CompressGZIP(file_openfga_v1alpha1_openfga_proto_rawDescData)
	})
	return file_openfga_v1alpha1_openfga_proto_rawDescData
}

var file_openfga_v1alpha1_openfga_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_openfga_v1alpha1_openfga_proto_goTypes = []interface{}{
	(*Module)(nil),                      // 0: fga.v1alpha1.Module
	(*Type)(nil),                        // 1: fga.v1alpha1.Type
	(*Relation)(nil),                    // 2: fga.v1alpha1.Relation
	(*Access)(nil),                      // 3: fga.v1alpha1.Access
	(*descriptorpb.ServiceOptions)(nil), // 4: google.protobuf.ServiceOptions
	(*descriptorpb.MethodOptions)(nil),  // 5: google.protobuf.MethodOptions
}
var file_openfga_v1alpha1_openfga_proto_depIdxs = []int32{
	1, // 0: fga.v1alpha1.Module.extends:type_name -> fga.v1alpha1.Type
	1, // 1: fga.v1alpha1.Module.types:type_name -> fga.v1alpha1.Type
	2, // 2: fga.v1alpha1.Type.relations:type_name -> fga.v1alpha1.Relation
	4, // 3: fga.v1alpha1.module:extendee -> google.protobuf.ServiceOptions
	5, // 4: fga.v1alpha1.access:extendee -> google.protobuf.MethodOptions
	0, // 5: fga.v1alpha1.module:type_name -> fga.v1alpha1.Module
	3, // 6: fga.v1alpha1.access:type_name -> fga.v1alpha1.Access
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	5, // [5:7] is the sub-list for extension type_name
	3, // [3:5] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_openfga_v1alpha1_openfga_proto_init() }
func file_openfga_v1alpha1_openfga_proto_init() {
	if File_openfga_v1alpha1_openfga_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_openfga_v1alpha1_openfga_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Module); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_openfga_v1alpha1_openfga_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Type); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_openfga_v1alpha1_openfga_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Relation); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_openfga_v1alpha1_openfga_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Access); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_openfga_v1alpha1_openfga_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_openfga_v1alpha1_openfga_proto_goTypes,
		DependencyIndexes: file_openfga_v1alpha1_openfga_proto_depIdxs,
		MessageInfos:      file_openfga_v1alpha1_openfga_proto_msgTypes,
		ExtensionInfos:    file_openfga_v1alpha1_openfga_proto_extTypes,
	}.Build()
	File_openfga_v1alpha1_openfga_proto = out.File
	file_openfga_v1alpha1_openfga_proto_rawDesc = nil
	file_openfga_v1alpha1_openfga_proto_goTypes = nil
	file_openfga_v1alpha1_openfga_proto_depIdxs = nil
}
