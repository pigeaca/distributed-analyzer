// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: proto/scheduler/scheduler.proto

package scheduler

import (
	task "github.com/distributedmarketplace/pkg/proto/task"
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

// SchedulerTask represents a task from the scheduler's perspective
type SchedulerTask struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        string                 `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Status        task.Status            `protobuf:"varint,2,opt,name=status,proto3,enum=task.Status" json:"status,omitempty"`
	SubTasks      []*task.SubTask        `protobuf:"bytes,3,rep,name=sub_tasks,json=subTasks,proto3" json:"sub_tasks,omitempty"`
	WorkerIds     []string               `protobuf:"bytes,4,rep,name=worker_ids,json=workerIds,proto3" json:"worker_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SchedulerTask) Reset() {
	*x = SchedulerTask{}
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SchedulerTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchedulerTask) ProtoMessage() {}

func (x *SchedulerTask) ProtoReflect() protoreflect.Message {
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchedulerTask.ProtoReflect.Descriptor instead.
func (*SchedulerTask) Descriptor() ([]byte, []int) {
	return file_proto_scheduler_scheduler_proto_rawDescGZIP(), []int{0}
}

func (x *SchedulerTask) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *SchedulerTask) GetStatus() task.Status {
	if x != nil {
		return x.Status
	}
	return task.Status(0)
}

func (x *SchedulerTask) GetSubTasks() []*task.SubTask {
	if x != nil {
		return x.SubTasks
	}
	return nil
}

func (x *SchedulerTask) GetWorkerIds() []string {
	if x != nil {
		return x.WorkerIds
	}
	return nil
}

// ScheduleTaskRequest is the request for scheduling a task
type ScheduleTaskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        string                 `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScheduleTaskRequest) Reset() {
	*x = ScheduleTaskRequest{}
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScheduleTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleTaskRequest) ProtoMessage() {}

func (x *ScheduleTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleTaskRequest.ProtoReflect.Descriptor instead.
func (*ScheduleTaskRequest) Descriptor() ([]byte, []int) {
	return file_proto_scheduler_scheduler_proto_rawDescGZIP(), []int{1}
}

func (x *ScheduleTaskRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

// ScheduleTaskResponse is the response for scheduling a task
type ScheduleTaskResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScheduleTaskResponse) Reset() {
	*x = ScheduleTaskResponse{}
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScheduleTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleTaskResponse) ProtoMessage() {}

func (x *ScheduleTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleTaskResponse.ProtoReflect.Descriptor instead.
func (*ScheduleTaskResponse) Descriptor() ([]byte, []int) {
	return file_proto_scheduler_scheduler_proto_rawDescGZIP(), []int{2}
}

func (x *ScheduleTaskResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// DivideTaskRequest is the request for dividing a task
type DivideTaskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        string                 `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DivideTaskRequest) Reset() {
	*x = DivideTaskRequest{}
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DivideTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DivideTaskRequest) ProtoMessage() {}

func (x *DivideTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DivideTaskRequest.ProtoReflect.Descriptor instead.
func (*DivideTaskRequest) Descriptor() ([]byte, []int) {
	return file_proto_scheduler_scheduler_proto_rawDescGZIP(), []int{3}
}

func (x *DivideTaskRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

// DivideTaskResponse is the response for dividing a task
type DivideTaskResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SubTasks      []*task.SubTask        `protobuf:"bytes,1,rep,name=sub_tasks,json=subTasks,proto3" json:"sub_tasks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DivideTaskResponse) Reset() {
	*x = DivideTaskResponse{}
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DivideTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DivideTaskResponse) ProtoMessage() {}

func (x *DivideTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DivideTaskResponse.ProtoReflect.Descriptor instead.
func (*DivideTaskResponse) Descriptor() ([]byte, []int) {
	return file_proto_scheduler_scheduler_proto_rawDescGZIP(), []int{4}
}

func (x *DivideTaskResponse) GetSubTasks() []*task.SubTask {
	if x != nil {
		return x.SubTasks
	}
	return nil
}

// AssignTaskRequest is the request for assigning a task
type AssignTaskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        string                 `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	WorkerId      string                 `protobuf:"bytes,2,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AssignTaskRequest) Reset() {
	*x = AssignTaskRequest{}
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AssignTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignTaskRequest) ProtoMessage() {}

func (x *AssignTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignTaskRequest.ProtoReflect.Descriptor instead.
func (*AssignTaskRequest) Descriptor() ([]byte, []int) {
	return file_proto_scheduler_scheduler_proto_rawDescGZIP(), []int{5}
}

func (x *AssignTaskRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *AssignTaskRequest) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

// AssignTaskResponse is the response for assigning a task
type AssignTaskResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AssignTaskResponse) Reset() {
	*x = AssignTaskResponse{}
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AssignTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignTaskResponse) ProtoMessage() {}

func (x *AssignTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignTaskResponse.ProtoReflect.Descriptor instead.
func (*AssignTaskResponse) Descriptor() ([]byte, []int) {
	return file_proto_scheduler_scheduler_proto_rawDescGZIP(), []int{6}
}

func (x *AssignTaskResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// GetTaskStatusRequest is the request for getting a task status
type GetTaskStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        string                 `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTaskStatusRequest) Reset() {
	*x = GetTaskStatusRequest{}
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTaskStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTaskStatusRequest) ProtoMessage() {}

func (x *GetTaskStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTaskStatusRequest.ProtoReflect.Descriptor instead.
func (*GetTaskStatusRequest) Descriptor() ([]byte, []int) {
	return file_proto_scheduler_scheduler_proto_rawDescGZIP(), []int{7}
}

func (x *GetTaskStatusRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

// GetTaskStatusResponse is the response for getting a task status
type GetTaskStatusResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        task.Status            `protobuf:"varint,1,opt,name=status,proto3,enum=task.Status" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTaskStatusResponse) Reset() {
	*x = GetTaskStatusResponse{}
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTaskStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTaskStatusResponse) ProtoMessage() {}

func (x *GetTaskStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_scheduler_scheduler_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTaskStatusResponse.ProtoReflect.Descriptor instead.
func (*GetTaskStatusResponse) Descriptor() ([]byte, []int) {
	return file_proto_scheduler_scheduler_proto_rawDescGZIP(), []int{8}
}

func (x *GetTaskStatusResponse) GetStatus() task.Status {
	if x != nil {
		return x.Status
	}
	return task.Status(0)
}

var File_proto_scheduler_scheduler_proto protoreflect.FileDescriptor

const file_proto_scheduler_scheduler_proto_rawDesc = "" +
	"\n" +
	"\x1fproto/scheduler/scheduler.proto\x12\tscheduler\x1a\x15proto/task/task.proto\"\x99\x01\n" +
	"\rSchedulerTask\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\tR\x06taskId\x12$\n" +
	"\x06status\x18\x02 \x01(\x0e2\f.task.StatusR\x06status\x12*\n" +
	"\tsub_tasks\x18\x03 \x03(\v2\r.task.SubTaskR\bsubTasks\x12\x1d\n" +
	"\n" +
	"worker_ids\x18\x04 \x03(\tR\tworkerIds\".\n" +
	"\x13ScheduleTaskRequest\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\tR\x06taskId\"0\n" +
	"\x14ScheduleTaskResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\",\n" +
	"\x11DivideTaskRequest\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\tR\x06taskId\"@\n" +
	"\x12DivideTaskResponse\x12*\n" +
	"\tsub_tasks\x18\x01 \x03(\v2\r.task.SubTaskR\bsubTasks\"I\n" +
	"\x11AssignTaskRequest\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\tR\x06taskId\x12\x1b\n" +
	"\tworker_id\x18\x02 \x01(\tR\bworkerId\".\n" +
	"\x12AssignTaskResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\"/\n" +
	"\x14GetTaskStatusRequest\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\tR\x06taskId\"=\n" +
	"\x15GetTaskStatusResponse\x12$\n" +
	"\x06status\x18\x01 \x01(\x0e2\f.task.StatusR\x06status2\xcd\x02\n" +
	"\x10SchedulerService\x12O\n" +
	"\fScheduleTask\x12\x1e.scheduler.ScheduleTaskRequest\x1a\x1f.scheduler.ScheduleTaskResponse\x12I\n" +
	"\n" +
	"DivideTask\x12\x1c.scheduler.DivideTaskRequest\x1a\x1d.scheduler.DivideTaskResponse\x12I\n" +
	"\n" +
	"AssignTask\x12\x1c.scheduler.AssignTaskRequest\x1a\x1d.scheduler.AssignTaskResponse\x12R\n" +
	"\rGetTaskStatus\x12\x1f.scheduler.GetTaskStatusRequest\x1a .scheduler.GetTaskStatusResponseB7Z5github.com/distributedmarketplace/pkg/proto/schedulerb\x06proto3"

var (
	file_proto_scheduler_scheduler_proto_rawDescOnce sync.Once
	file_proto_scheduler_scheduler_proto_rawDescData []byte
)

func file_proto_scheduler_scheduler_proto_rawDescGZIP() []byte {
	file_proto_scheduler_scheduler_proto_rawDescOnce.Do(func() {
		file_proto_scheduler_scheduler_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_scheduler_scheduler_proto_rawDesc), len(file_proto_scheduler_scheduler_proto_rawDesc)))
	})
	return file_proto_scheduler_scheduler_proto_rawDescData
}

var file_proto_scheduler_scheduler_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_scheduler_scheduler_proto_goTypes = []any{
	(*SchedulerTask)(nil),         // 0: scheduler.SchedulerTask
	(*ScheduleTaskRequest)(nil),   // 1: scheduler.ScheduleTaskRequest
	(*ScheduleTaskResponse)(nil),  // 2: scheduler.ScheduleTaskResponse
	(*DivideTaskRequest)(nil),     // 3: scheduler.DivideTaskRequest
	(*DivideTaskResponse)(nil),    // 4: scheduler.DivideTaskResponse
	(*AssignTaskRequest)(nil),     // 5: scheduler.AssignTaskRequest
	(*AssignTaskResponse)(nil),    // 6: scheduler.AssignTaskResponse
	(*GetTaskStatusRequest)(nil),  // 7: scheduler.GetTaskStatusRequest
	(*GetTaskStatusResponse)(nil), // 8: scheduler.GetTaskStatusResponse
	(task.Status)(0),              // 9: task.Status
	(*task.SubTask)(nil),          // 10: task.SubTask
}
var file_proto_scheduler_scheduler_proto_depIdxs = []int32{
	9,  // 0: scheduler.SchedulerTask.status:type_name -> task.Status
	10, // 1: scheduler.SchedulerTask.sub_tasks:type_name -> task.SubTask
	10, // 2: scheduler.DivideTaskResponse.sub_tasks:type_name -> task.SubTask
	9,  // 3: scheduler.GetTaskStatusResponse.status:type_name -> task.Status
	1,  // 4: scheduler.SchedulerService.ScheduleTask:input_type -> scheduler.ScheduleTaskRequest
	3,  // 5: scheduler.SchedulerService.DivideTask:input_type -> scheduler.DivideTaskRequest
	5,  // 6: scheduler.SchedulerService.AssignTask:input_type -> scheduler.AssignTaskRequest
	7,  // 7: scheduler.SchedulerService.GetTaskStatus:input_type -> scheduler.GetTaskStatusRequest
	2,  // 8: scheduler.SchedulerService.ScheduleTask:output_type -> scheduler.ScheduleTaskResponse
	4,  // 9: scheduler.SchedulerService.DivideTask:output_type -> scheduler.DivideTaskResponse
	6,  // 10: scheduler.SchedulerService.AssignTask:output_type -> scheduler.AssignTaskResponse
	8,  // 11: scheduler.SchedulerService.GetTaskStatus:output_type -> scheduler.GetTaskStatusResponse
	8,  // [8:12] is the sub-list for method output_type
	4,  // [4:8] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_proto_scheduler_scheduler_proto_init() }
func file_proto_scheduler_scheduler_proto_init() {
	if File_proto_scheduler_scheduler_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_scheduler_scheduler_proto_rawDesc), len(file_proto_scheduler_scheduler_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_scheduler_scheduler_proto_goTypes,
		DependencyIndexes: file_proto_scheduler_scheduler_proto_depIdxs,
		MessageInfos:      file_proto_scheduler_scheduler_proto_msgTypes,
	}.Build()
	File_proto_scheduler_scheduler_proto = out.File
	file_proto_scheduler_scheduler_proto_goTypes = nil
	file_proto_scheduler_scheduler_proto_depIdxs = nil
}
