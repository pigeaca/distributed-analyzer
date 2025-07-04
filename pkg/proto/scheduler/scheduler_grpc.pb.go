// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.29.3
// source: proto/scheduler/scheduler.proto

package scheduler

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SchedulerService_ScheduleTask_FullMethodName  = "/scheduler.SchedulerService/ScheduleTask"
	SchedulerService_DivideTask_FullMethodName    = "/scheduler.SchedulerService/DivideTask"
	SchedulerService_AssignTask_FullMethodName    = "/scheduler.SchedulerService/AssignTask"
	SchedulerService_GetTaskStatus_FullMethodName = "/scheduler.SchedulerService/GetTaskStatus"
)

// SchedulerServiceClient is the client API for SchedulerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SchedulerServiceClient interface {
	// ScheduleTask assigns a task to appropriate workers
	ScheduleTask(ctx context.Context, in *ScheduleTaskRequest, opts ...grpc.CallOption) (*ScheduleTaskResponse, error)
	// DivideTask splits a task into subtasks if needed
	DivideTask(ctx context.Context, in *DivideTaskRequest, opts ...grpc.CallOption) (*DivideTaskResponse, error)
	// AssignTask assigns a task or subtask to a specific worker
	AssignTask(ctx context.Context, in *AssignTaskRequest, opts ...grpc.CallOption) (*AssignTaskResponse, error)
	// GetTaskStatus retrieves the current status of a task
	GetTaskStatus(ctx context.Context, in *GetTaskStatusRequest, opts ...grpc.CallOption) (*GetTaskStatusResponse, error)
}

type schedulerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSchedulerServiceClient(cc grpc.ClientConnInterface) SchedulerServiceClient {
	return &schedulerServiceClient{cc}
}

func (c *schedulerServiceClient) ScheduleTask(ctx context.Context, in *ScheduleTaskRequest, opts ...grpc.CallOption) (*ScheduleTaskResponse, error) {
	out := new(ScheduleTaskResponse)
	err := c.cc.Invoke(ctx, SchedulerService_ScheduleTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerServiceClient) DivideTask(ctx context.Context, in *DivideTaskRequest, opts ...grpc.CallOption) (*DivideTaskResponse, error) {
	out := new(DivideTaskResponse)
	err := c.cc.Invoke(ctx, SchedulerService_DivideTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerServiceClient) AssignTask(ctx context.Context, in *AssignTaskRequest, opts ...grpc.CallOption) (*AssignTaskResponse, error) {
	out := new(AssignTaskResponse)
	err := c.cc.Invoke(ctx, SchedulerService_AssignTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerServiceClient) GetTaskStatus(ctx context.Context, in *GetTaskStatusRequest, opts ...grpc.CallOption) (*GetTaskStatusResponse, error) {
	out := new(GetTaskStatusResponse)
	err := c.cc.Invoke(ctx, SchedulerService_GetTaskStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SchedulerServiceServer is the server API for SchedulerService service.
// All implementations must embed UnimplementedSchedulerServiceServer
// for forward compatibility
type SchedulerServiceServer interface {
	// ScheduleTask assigns a task to appropriate workers
	ScheduleTask(context.Context, *ScheduleTaskRequest) (*ScheduleTaskResponse, error)
	// DivideTask splits a task into subtasks if needed
	DivideTask(context.Context, *DivideTaskRequest) (*DivideTaskResponse, error)
	// AssignTask assigns a task or subtask to a specific worker
	AssignTask(context.Context, *AssignTaskRequest) (*AssignTaskResponse, error)
	// GetTaskStatus retrieves the current status of a task
	GetTaskStatus(context.Context, *GetTaskStatusRequest) (*GetTaskStatusResponse, error)
	mustEmbedUnimplementedSchedulerServiceServer()
}

// UnimplementedSchedulerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSchedulerServiceServer struct {
}

func (UnimplementedSchedulerServiceServer) ScheduleTask(context.Context, *ScheduleTaskRequest) (*ScheduleTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleTask not implemented")
}
func (UnimplementedSchedulerServiceServer) DivideTask(context.Context, *DivideTaskRequest) (*DivideTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DivideTask not implemented")
}
func (UnimplementedSchedulerServiceServer) AssignTask(context.Context, *AssignTaskRequest) (*AssignTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignTask not implemented")
}
func (UnimplementedSchedulerServiceServer) GetTaskStatus(context.Context, *GetTaskStatusRequest) (*GetTaskStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskStatus not implemented")
}
func (UnimplementedSchedulerServiceServer) mustEmbedUnimplementedSchedulerServiceServer() {}

// UnsafeSchedulerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SchedulerServiceServer will
// result in compilation errors.
type UnsafeSchedulerServiceServer interface {
	mustEmbedUnimplementedSchedulerServiceServer()
}

func RegisterSchedulerServiceServer(s grpc.ServiceRegistrar, srv SchedulerServiceServer) {
	s.RegisterService(&SchedulerService_ServiceDesc, srv)
}

func _SchedulerService_ScheduleTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScheduleTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServiceServer).ScheduleTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchedulerService_ScheduleTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServiceServer).ScheduleTask(ctx, req.(*ScheduleTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchedulerService_DivideTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DivideTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServiceServer).DivideTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchedulerService_DivideTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServiceServer).DivideTask(ctx, req.(*DivideTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchedulerService_AssignTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServiceServer).AssignTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchedulerService_AssignTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServiceServer).AssignTask(ctx, req.(*AssignTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchedulerService_GetTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServiceServer).GetTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchedulerService_GetTaskStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServiceServer).GetTaskStatus(ctx, req.(*GetTaskStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SchedulerService_ServiceDesc is the grpc.ServiceDesc for SchedulerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SchedulerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "scheduler.SchedulerService",
	HandlerType: (*SchedulerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ScheduleTask",
			Handler:    _SchedulerService_ScheduleTask_Handler,
		},
		{
			MethodName: "DivideTask",
			Handler:    _SchedulerService_DivideTask_Handler,
		},
		{
			MethodName: "AssignTask",
			Handler:    _SchedulerService_AssignTask_Handler,
		},
		{
			MethodName: "GetTaskStatus",
			Handler:    _SchedulerService_GetTaskStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/scheduler/scheduler.proto",
}
