package descrule

import (
	"reflect"
	"testing"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestAllDescriptors(t *testing.T) {
	f := builder.NewFile("test_proto").AddMessage(
		builder.NewMessage("Message1").AddField(builder.NewField("a", builder.FieldTypeString())),
	).AddMessage(
		builder.NewMessage("Message2").AddField(
			builder.NewField("a", builder.FieldTypeString()),
		).AddField(
			builder.NewField("b", builder.FieldTypeString()),
		).AddNestedMessage(
			builder.NewMessage("NestMessage1").AddField(
				builder.NewField("a", builder.FieldTypeString()),
			).AddNestedEnum(
				builder.NewEnum("Enum").AddValue(
					builder.NewEnumValue("EnumValue1"),
				),
			),
		),
	).AddService(
		builder.NewService("Service1").AddMethod(
			builder.NewMethod("Rpc1",
				builder.RpcTypeMessage(builder.NewMessage("Request"), false),
				builder.RpcTypeMessage(builder.NewMessage("Response"), false),
			),
		),
	)
	fd, err := f.Build()
	if err != nil {
		t.Fatal(err)
	}
	descriptors := AllDescriptors(fd)
	// Total: files + messages + fields + services + methods + enums + enum_values
	wantNum := 1 + 3 + 4 + 1 + 1 + 1 + 1
	gotNum := len(descriptors)
	if gotNum != wantNum {
		t.Errorf("AllDescriptors got %d descriptors, but want %d", gotNum, wantNum)
	}
}

func TestSourceInfo(t *testing.T) {
	f := builder.NewFile("test_proto").AddMessage(
		builder.NewMessage("Message1").AddField(
			builder.NewField("a", builder.FieldTypeString()).SetComments(
				builder.Comments{
					LeadingComment: "Message1_Field1_Comments",
				},
			),
		).AddNestedEnum(
			builder.NewEnum("Enum").AddValue(
				builder.NewEnumValue("EnumValue1").SetComments(
					builder.Comments{
						LeadingComment: "Message1_Enum1_Value1_Comments",
					}),
			).SetComments(
				builder.Comments{
					LeadingComment: "Message1_Enum1_Comments",
				},
			),
		).SetComments(
			builder.Comments{
				LeadingComment: "Message1_Comments",
			}),
	).AddService(
		builder.NewService("Service1").AddMethod(
			builder.NewMethod("Rpc1",
				builder.RpcTypeMessage(builder.NewMessage("Request"), false),
				builder.RpcTypeMessage(builder.NewMessage("Response"), false),
			).SetComments(
				builder.Comments{
					LeadingComment: "Service1_Method1_Comments",
				},
			),
		).SetComments(
			builder.Comments{
				LeadingComment: "Service1_Comments",
			},
		),
	)
	fd, err := f.Build()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name       string
		descriptor lint.Descriptor
		sourceInfo lint.SourceInfo
	}{
		{"File", NewFile(fd), &sourceInfo{file: &fileInfo{path: "test_proto"}}},
		{"Message", NewMessage(fd.GetMessageTypes()[0]), &sourceInfo{file: &fileInfo{path: "test_proto"}, leadingComments: "Message1_Comments"}},
		{"Field", NewField(fd.GetMessageTypes()[0].GetFields()[0]), &sourceInfo{file: &fileInfo{path: "test_proto"}, leadingComments: "Message1_Field1_Comments"}},
		{"Enum", NewEnum(fd.GetMessageTypes()[0].GetNestedEnumTypes()[0]), &sourceInfo{file: &fileInfo{path: "test_proto"}, leadingComments: "Message1_Enum1_Comments"}},
		{"EnumValue", NewEnumValue(fd.GetMessageTypes()[0].GetNestedEnumTypes()[0].GetValues()[0]), &sourceInfo{file: &fileInfo{path: "test_proto"}, leadingComments: "Message1_Enum1_Value1_Comments"}},
		{"Service", NewService(fd.GetServices()[0]), &sourceInfo{file: &fileInfo{path: "test_proto"}, leadingComments: "Service1_Comments"}},
		{"Method", NewMethod(fd.GetServices()[0].GetMethods()[0]), &sourceInfo{file: &fileInfo{path: "test_proto"}, leadingComments: "Service1_Method1_Comments"}},
		{"Descriptor", NewDescriptor(fd.GetServices()[0].GetMethods()[0]), &sourceInfo{file: &fileInfo{path: "test_proto"}, leadingComments: "Service1_Method1_Comments"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.descriptor.SourceInfo()
			if !reflect.DeepEqual(got, test.sourceInfo) {
				t.Errorf("SourceInfo returns %v, but want %v", got, test.sourceInfo)
			}
		})
	}
}
