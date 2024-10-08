// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: petstore.proto

package petstorev1

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
	PetstoreService_GetPetByID_FullMethodName        = "/petstore.v1.PetstoreService/GetPetByID"
	PetstoreService_GetPets_FullMethodName           = "/petstore.v1.PetstoreService/GetPets"
	PetstoreService_UpdatePetWithForm_FullMethodName = "/petstore.v1.PetstoreService/UpdatePetWithForm"
	PetstoreService_DeletePet_FullMethodName         = "/petstore.v1.PetstoreService/DeletePet"
	PetstoreService_UploadFile_FullMethodName        = "/petstore.v1.PetstoreService/UploadFile"
	PetstoreService_AddPet_FullMethodName            = "/petstore.v1.PetstoreService/AddPet"
	PetstoreService_UpdatePet_FullMethodName         = "/petstore.v1.PetstoreService/UpdatePet"
	PetstoreService_FindPetsByStatus_FullMethodName  = "/petstore.v1.PetstoreService/FindPetsByStatus"
)

// PetstoreServiceClient is the client API for PetstoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PetstoreServiceClient interface {
	GetPetByID(ctx context.Context, in *PetID, opts ...grpc.CallOption) (*Pet, error)
	GetPets(ctx context.Context, in *Empty, opts ...grpc.CallOption) (PetstoreService_GetPetsClient, error)
	UpdatePetWithForm(ctx context.Context, in *UpdatePetWithFormReq, opts ...grpc.CallOption) (*Empty, error)
	DeletePet(ctx context.Context, in *PetID, opts ...grpc.CallOption) (*Empty, error)
	UploadFile(ctx context.Context, in *UploadFileReq, opts ...grpc.CallOption) (*ApiResponse, error)
	AddPet(ctx context.Context, in *Pet, opts ...grpc.CallOption) (*Pet, error)
	UpdatePet(ctx context.Context, in *Pet, opts ...grpc.CallOption) (*Pet, error)
	FindPetsByStatus(ctx context.Context, in *StatusReq, opts ...grpc.CallOption) (*Pets, error)
}

type petstoreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPetstoreServiceClient(cc grpc.ClientConnInterface) PetstoreServiceClient {
	return &petstoreServiceClient{cc}
}

func (c *petstoreServiceClient) GetPetByID(ctx context.Context, in *PetID, opts ...grpc.CallOption) (*Pet, error) {
	out := new(Pet)
	err := c.cc.Invoke(ctx, PetstoreService_GetPetByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petstoreServiceClient) GetPets(ctx context.Context, in *Empty, opts ...grpc.CallOption) (PetstoreService_GetPetsClient, error) {
	stream, err := c.cc.NewStream(ctx, &PetstoreService_ServiceDesc.Streams[0], PetstoreService_GetPets_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &petstoreServiceGetPetsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PetstoreService_GetPetsClient interface {
	Recv() (*Pet, error)
	grpc.ClientStream
}

type petstoreServiceGetPetsClient struct {
	grpc.ClientStream
}

func (x *petstoreServiceGetPetsClient) Recv() (*Pet, error) {
	m := new(Pet)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *petstoreServiceClient) UpdatePetWithForm(ctx context.Context, in *UpdatePetWithFormReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, PetstoreService_UpdatePetWithForm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petstoreServiceClient) DeletePet(ctx context.Context, in *PetID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, PetstoreService_DeletePet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petstoreServiceClient) UploadFile(ctx context.Context, in *UploadFileReq, opts ...grpc.CallOption) (*ApiResponse, error) {
	out := new(ApiResponse)
	err := c.cc.Invoke(ctx, PetstoreService_UploadFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petstoreServiceClient) AddPet(ctx context.Context, in *Pet, opts ...grpc.CallOption) (*Pet, error) {
	out := new(Pet)
	err := c.cc.Invoke(ctx, PetstoreService_AddPet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petstoreServiceClient) UpdatePet(ctx context.Context, in *Pet, opts ...grpc.CallOption) (*Pet, error) {
	out := new(Pet)
	err := c.cc.Invoke(ctx, PetstoreService_UpdatePet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petstoreServiceClient) FindPetsByStatus(ctx context.Context, in *StatusReq, opts ...grpc.CallOption) (*Pets, error) {
	out := new(Pets)
	err := c.cc.Invoke(ctx, PetstoreService_FindPetsByStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PetstoreServiceServer is the server API for PetstoreService service.
// All implementations must embed UnimplementedPetstoreServiceServer
// for forward compatibility
type PetstoreServiceServer interface {
	GetPetByID(context.Context, *PetID) (*Pet, error)
	GetPets(*Empty, PetstoreService_GetPetsServer) error
	UpdatePetWithForm(context.Context, *UpdatePetWithFormReq) (*Empty, error)
	DeletePet(context.Context, *PetID) (*Empty, error)
	UploadFile(context.Context, *UploadFileReq) (*ApiResponse, error)
	AddPet(context.Context, *Pet) (*Pet, error)
	UpdatePet(context.Context, *Pet) (*Pet, error)
	FindPetsByStatus(context.Context, *StatusReq) (*Pets, error)
	mustEmbedUnimplementedPetstoreServiceServer()
}

// UnimplementedPetstoreServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPetstoreServiceServer struct {
}

func (UnimplementedPetstoreServiceServer) GetPetByID(context.Context, *PetID) (*Pet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPetByID not implemented")
}
func (UnimplementedPetstoreServiceServer) GetPets(*Empty, PetstoreService_GetPetsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetPets not implemented")
}
func (UnimplementedPetstoreServiceServer) UpdatePetWithForm(context.Context, *UpdatePetWithFormReq) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePetWithForm not implemented")
}
func (UnimplementedPetstoreServiceServer) DeletePet(context.Context, *PetID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePet not implemented")
}
func (UnimplementedPetstoreServiceServer) UploadFile(context.Context, *UploadFileReq) (*ApiResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedPetstoreServiceServer) AddPet(context.Context, *Pet) (*Pet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPet not implemented")
}
func (UnimplementedPetstoreServiceServer) UpdatePet(context.Context, *Pet) (*Pet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePet not implemented")
}
func (UnimplementedPetstoreServiceServer) FindPetsByStatus(context.Context, *StatusReq) (*Pets, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindPetsByStatus not implemented")
}
func (UnimplementedPetstoreServiceServer) mustEmbedUnimplementedPetstoreServiceServer() {}

// UnsafePetstoreServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PetstoreServiceServer will
// result in compilation errors.
type UnsafePetstoreServiceServer interface {
	mustEmbedUnimplementedPetstoreServiceServer()
}

func RegisterPetstoreServiceServer(s grpc.ServiceRegistrar, srv PetstoreServiceServer) {
	s.RegisterService(&PetstoreService_ServiceDesc, srv)
}

func _PetstoreService_GetPetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetstoreServiceServer).GetPetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetstoreService_GetPetByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetstoreServiceServer).GetPetByID(ctx, req.(*PetID))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetstoreService_GetPets_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PetstoreServiceServer).GetPets(m, &petstoreServiceGetPetsServer{stream})
}

type PetstoreService_GetPetsServer interface {
	Send(*Pet) error
	grpc.ServerStream
}

type petstoreServiceGetPetsServer struct {
	grpc.ServerStream
}

func (x *petstoreServiceGetPetsServer) Send(m *Pet) error {
	return x.ServerStream.SendMsg(m)
}

func _PetstoreService_UpdatePetWithForm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePetWithFormReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetstoreServiceServer).UpdatePetWithForm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetstoreService_UpdatePetWithForm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetstoreServiceServer).UpdatePetWithForm(ctx, req.(*UpdatePetWithFormReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetstoreService_DeletePet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetstoreServiceServer).DeletePet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetstoreService_DeletePet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetstoreServiceServer).DeletePet(ctx, req.(*PetID))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetstoreService_UploadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadFileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetstoreServiceServer).UploadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetstoreService_UploadFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetstoreServiceServer).UploadFile(ctx, req.(*UploadFileReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetstoreService_AddPet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetstoreServiceServer).AddPet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetstoreService_AddPet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetstoreServiceServer).AddPet(ctx, req.(*Pet))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetstoreService_UpdatePet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetstoreServiceServer).UpdatePet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetstoreService_UpdatePet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetstoreServiceServer).UpdatePet(ctx, req.(*Pet))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetstoreService_FindPetsByStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetstoreServiceServer).FindPetsByStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetstoreService_FindPetsByStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetstoreServiceServer).FindPetsByStatus(ctx, req.(*StatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

// PetstoreService_ServiceDesc is the grpc.ServiceDesc for PetstoreService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PetstoreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "petstore.v1.PetstoreService",
	HandlerType: (*PetstoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPetByID",
			Handler:    _PetstoreService_GetPetByID_Handler,
		},
		{
			MethodName: "UpdatePetWithForm",
			Handler:    _PetstoreService_UpdatePetWithForm_Handler,
		},
		{
			MethodName: "DeletePet",
			Handler:    _PetstoreService_DeletePet_Handler,
		},
		{
			MethodName: "UploadFile",
			Handler:    _PetstoreService_UploadFile_Handler,
		},
		{
			MethodName: "AddPet",
			Handler:    _PetstoreService_AddPet_Handler,
		},
		{
			MethodName: "UpdatePet",
			Handler:    _PetstoreService_UpdatePet_Handler,
		},
		{
			MethodName: "FindPetsByStatus",
			Handler:    _PetstoreService_FindPetsByStatus_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetPets",
			Handler:       _PetstoreService_GetPets_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "petstore.proto",
}
