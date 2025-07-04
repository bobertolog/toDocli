// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: proto/task.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Task struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Status        string                 `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Task) Reset() {
	*x = Task{}
	mi := &file_proto_task_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{0}
}

func (x *Task) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Task) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Task) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Task) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type TaskID struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskID) Reset() {
	*x = TaskID{}
	mi := &file_proto_task_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskID) ProtoMessage() {}

func (x *TaskID) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskID.ProtoReflect.Descriptor instead.
func (*TaskID) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{1}
}

func (x *TaskID) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type TaskList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Tasks         []*Task                `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskList) Reset() {
	*x = TaskList{}
	mi := &file_proto_task_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskList) ProtoMessage() {}

func (x *TaskList) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskList.ProtoReflect.Descriptor instead.
func (*TaskList) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{2}
}

func (x *TaskList) GetTasks() []*Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_proto_task_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_task_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_proto_task_proto_rawDescGZIP(), []int{3}
}

var File_proto_task_proto protoreflect.FileDescriptor

const file_proto_task_proto_rawDesc = "" +
	"\n" +
	"\x10proto/task.proto\x12\x04task\"f\n" +
	"\x04Task\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12\x16\n" +
	"\x06status\x18\x04 \x01(\tR\x06status\"\x18\n" +
	"\x06TaskID\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\",\n" +
	"\bTaskList\x12 \n" +
	"\x05tasks\x18\x01 \x03(\v2\n" +
	".task.TaskR\x05tasks\"\a\n" +
	"\x05Empty2\xd3\x01\n" +
	"\vTaskService\x12&\n" +
	"\n" +
	"CreateTask\x12\n" +
	".task.Task\x1a\f.task.TaskID\x12#\n" +
	"\aGetTask\x12\f.task.TaskID\x1a\n" +
	".task.Task\x12(\n" +
	"\tListTasks\x12\v.task.Empty\x1a\x0e.task.TaskList\x12$\n" +
	"\n" +
	"UpdateTask\x12\n" +
	".task.Task\x1a\n" +
	".task.Task\x12'\n" +
	"\n" +
	"DeleteTask\x12\f.task.TaskID\x1a\v.task.EmptyB\fZ\n" +
	"todocli/pbb\x06proto3"

var (
	file_proto_task_proto_rawDescOnce sync.Once
	file_proto_task_proto_rawDescData []byte
)

func file_proto_task_proto_rawDescGZIP() []byte {
	file_proto_task_proto_rawDescOnce.Do(func() {
		file_proto_task_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_task_proto_rawDesc), len(file_proto_task_proto_rawDesc)))
	})
	return file_proto_task_proto_rawDescData
}

var file_proto_task_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_task_proto_goTypes = []any{
	(*Task)(nil),     // 0: task.Task
	(*TaskID)(nil),   // 1: task.TaskID
	(*TaskList)(nil), // 2: task.TaskList
	(*Empty)(nil),    // 3: task.Empty
}
var file_proto_task_proto_depIdxs = []int32{
	0, // 0: task.TaskList.tasks:type_name -> task.Task
	0, // 1: task.TaskService.CreateTask:input_type -> task.Task
	1, // 2: task.TaskService.GetTask:input_type -> task.TaskID
	3, // 3: task.TaskService.ListTasks:input_type -> task.Empty
	0, // 4: task.TaskService.UpdateTask:input_type -> task.Task
	1, // 5: task.TaskService.DeleteTask:input_type -> task.TaskID
	1, // 6: task.TaskService.CreateTask:output_type -> task.TaskID
	0, // 7: task.TaskService.GetTask:output_type -> task.Task
	2, // 8: task.TaskService.ListTasks:output_type -> task.TaskList
	0, // 9: task.TaskService.UpdateTask:output_type -> task.Task
	3, // 10: task.TaskService.DeleteTask:output_type -> task.Empty
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_task_proto_init() }
func file_proto_task_proto_init() {
	if File_proto_task_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_task_proto_rawDesc), len(file_proto_task_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_task_proto_goTypes,
		DependencyIndexes: file_proto_task_proto_depIdxs,
		MessageInfos:      file_proto_task_proto_msgTypes,
	}.Build()
	File_proto_task_proto = out.File
	file_proto_task_proto_goTypes = nil
	file_proto_task_proto_depIdxs = nil
}
