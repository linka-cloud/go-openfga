package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	resource "go.linka.cloud/go-openfga/example/pb"
)

func TestDeclarationOrder(t *testing.T) {
	file := resource.File_example_pb_resource_proto
	req := &pluginpb.CodeGeneratorRequest{FileToGenerate: []string{file.Path()}}
	seen := make(map[string]bool)
	var addFile func(protoreflect.FileDescriptor)
	addFile = func(f protoreflect.FileDescriptor) {
		if seen[f.Path()] {
			return
		}
		seen[f.Path()] = true
		for i := 0; i < f.Imports().Len(); i++ {
			addFile(f.Imports().Get(i).FileDescriptor)
		}
		req.ProtoFile = append(req.ProtoFile, protodesc.ToFileDescriptorProto(f))
	}
	addFile(file)
	// Compiled descriptors omit source locations required by the template.
	req.ProtoFile[len(req.ProtoFile)-1].SourceCodeInfo = &descriptorpb.SourceCodeInfo{
		Location: []*descriptorpb.SourceCodeInfo_Location{{Path: []int32{12}, Span: []int32{0, 0, 0}}},
	}
	input, err := proto.Marshal(req)
	require.NoError(t, err)

	var first string
	for i := range 32 {
		var output bytes.Buffer
		pgs.Init(pgs.ProtocInput(bytes.NewReader(input)), pgs.ProtocOutput(&output)).
			RegisterModule(FGA()).
			RegisterPostProcessor(pgsgo.GoFmt()).
			Render()
		var resp pluginpb.CodeGeneratorResponse
		require.NoError(t, proto.Unmarshal(output.Bytes(), &resp))
		require.Empty(t, resp.GetError())
		require.Len(t, resp.File, 2)
		var content string
		for _, f := range resp.File {
			if f.GetName() == "resource/resource.pb.fga.go" {
				content = f.GetContent()
			}
		}
		require.NotEmpty(t, content)
		if i == 0 {
			first = content
		} else {
			require.Equal(t, first, content, "generation %d changed declaration order", i)
		}
	}

	generated, err := parser.ParseFile(token.NewFileSet(), "resource.fga.go", first, 0)
	require.NoError(t, err)
	var constants []string
	for _, decl := range generated.Decls {
		if d, ok := decl.(*ast.GenDecl); ok && d.Tok == token.CONST {
			for _, spec := range d.Specs {
				constants = append(constants, spec.(*ast.ValueSpec).Names[0].Name)
			}
		}
	}
	require.Equal(t, []string{
		"FGAResourceType", "FGAResourceAdmin", "FGAResourceCanAddSub", "FGAResourceCanDelete",
		"FGAResourceCanRead", "FGAResourceCanReadSub", "FGAResourceCanUpdate", "FGAResourceReader", "FGAResourceSystem",
		"FGASubType", "FGASubAdmin", "FGASubCanRead", "FGASubReader", "FGASubResource",
		"FGASystemType", "FGASystemCanResourceCreate", "FGASystemCanResourceDelete", "FGASystemCanResourceList",
		"FGASystemCanResourceRead", "FGASystemCanResourceUpdate", "FGASystemCanResourceWatch",
		"FGASystemResourceAdmin", "FGASystemResourceReader", "FGASystemResourceWatcher", "FGASystemResourceWriter",
	}, constants)
}
