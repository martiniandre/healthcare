package pb

import (
	"context"
	"google.golang.org/grpc"
)

type CreateAuditLogRequest struct {
	CorrelationId string
	CallerUserId  string
	CallerRole    string
	Method        string
	AccessGranted bool
}

type CreateAuditLogResponse struct {
	Id string
}

type ListAuditLogsRequest struct {
	Limit  int32
	Offset int32
}

type AuditLog struct {
	Id            string
	CorrelationId string
	CallerUserId  string
	CallerRole    string
	Method        string
	AccessGranted bool
	CreatedAt     string
}

type ListAuditLogsResponse struct {
	AuditLogs []*AuditLog
	Total     int32
}

type AuditLogsServiceServer interface {
	CreateAuditLog(context.Context, *CreateAuditLogRequest) (*CreateAuditLogResponse, error)
	ListAuditLogs(context.Context, *ListAuditLogsRequest) (*ListAuditLogsResponse, error)
}

func RegisterAuditLogsServiceServer(registrar grpc.ServiceRegistrar, serviceImpl AuditLogsServiceServer) {
	registrar.RegisterService(&AuditLogsService_ServiceDesc, serviceImpl)
}

var AuditLogsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "audit_logs.v1.AuditLogsService",
	HandlerType: (*AuditLogsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAuditLog",
			Handler: func(serviceImpl interface{}, contextVal context.Context, decoder func(interface{}) error, unaryInterceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				request := new(CreateAuditLogRequest)
				if decodeError := decoder(request); decodeError != nil {
					return nil, decodeError
				}
				if unaryInterceptor == nil {
					return serviceImpl.(AuditLogsServiceServer).CreateAuditLog(contextVal, request)
				}
				info := &grpc.UnaryServerInfo{
					Server:     serviceImpl,
					FullMethod: "/audit_logs.v1.AuditLogsService/CreateAuditLog",
				}
				handler := func(innerContext context.Context, innerRequest interface{}) (interface{}, error) {
					return serviceImpl.(AuditLogsServiceServer).CreateAuditLog(innerContext, innerRequest.(*CreateAuditLogRequest))
				}
				return unaryInterceptor(contextVal, request, info, handler)
			},
		},
		{
			MethodName: "ListAuditLogs",
			Handler: func(serviceImpl interface{}, contextVal context.Context, decoder func(interface{}) error, unaryInterceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				request := new(ListAuditLogsRequest)
				if decodeError := decoder(request); decodeError != nil {
					return nil, decodeError
				}
				if unaryInterceptor == nil {
					return serviceImpl.(AuditLogsServiceServer).ListAuditLogs(contextVal, request)
				}
				info := &grpc.UnaryServerInfo{
					Server:     serviceImpl,
					FullMethod: "/audit_logs.v1.AuditLogsService/ListAuditLogs",
				}
				handler := func(innerContext context.Context, innerRequest interface{}) (interface{}, error) {
					return serviceImpl.(AuditLogsServiceServer).ListAuditLogs(innerContext, innerRequest.(*ListAuditLogsRequest))
				}
				return unaryInterceptor(contextVal, request, info, handler)
			},
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "audit_logs.proto",
}
