package pb

import (
	"context"
	"google.golang.org/grpc"
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token  string
	UserId string
	Role   string
}

type RegisterRequest struct {
	Email    string
	Password string
	FullName string
	Role     string
}

type RegisterResponse struct {
	UserId string
}

type LogoutRequest struct{}
type LogoutResponse struct{}

type AuthServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.v1.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(LoginRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.(AuthServiceServer).Login(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/auth.v1.AuthService/Login",
				}
				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return srv.(AuthServiceServer).Login(ctx, req.(*LoginRequest))
				}
				return interceptor(ctx, in, info, handler)
			},
		},
		{
			MethodName: "Register",
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(RegisterRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.(AuthServiceServer).Register(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/auth.v1.AuthService/Register",
				}
				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return srv.(AuthServiceServer).Register(ctx, req.(*RegisterRequest))
				}
				return interceptor(ctx, in, info, handler)
			},
		},
		{
			MethodName: "Logout",
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(LogoutRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.(AuthServiceServer).Logout(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/auth.v1.AuthService/Logout",
				}
				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return srv.(AuthServiceServer).Logout(ctx, req.(*LogoutRequest))
				}
				return interceptor(ctx, in, info, handler)
			},
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}
