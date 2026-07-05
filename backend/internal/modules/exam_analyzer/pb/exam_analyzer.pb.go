package pb

import (
	"context"
	"google.golang.org/grpc"
)

type ExamAnalysisItem struct {
	Id               string
	UserId           string
	PatientFhirId    string
	ExamType         string
	FileName         string
	Status           string
	AnalysisResponse string
	ConsentGiven     bool
	Anonymized       bool
	CreatedAt        string
	UpdatedAt        string
}

type ListAnalysesRequest struct {
	PatientFhirId string
}

type ListAnalysesResponse struct {
	Analyses []*ExamAnalysisItem
}

type GetAnalysisRequest struct {
	AnalysisId string
}

type GetAnalysisResponse struct {
	Analysis *ExamAnalysisItem
}

type DeleteAnalysisRequest struct {
	AnalysisId string
}

type DeleteAnalysisResponse struct {
	Success bool
}

type ExamAnalyzerServiceServer interface {
	ListAnalyses(ctx context.Context, req *ListAnalysesRequest) (*ListAnalysesResponse, error)
	GetAnalysis(ctx context.Context, req *GetAnalysisRequest) (*GetAnalysisResponse, error)
	DeleteAnalysis(ctx context.Context, req *DeleteAnalysisRequest) (*DeleteAnalysisResponse, error)
}

func RegisterExamAnalyzerServiceServer(serverRegistrar grpc.ServiceRegistrar, server ExamAnalyzerServiceServer) {
	serverRegistrar.RegisterService(&ExamAnalyzerService_ServiceDesc, server)
}

var ExamAnalyzerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "exam_analyzer.v1.ExamAnalyzerService",
	HandlerType: (*ExamAnalyzerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAnalyses",
			Handler: func(serverInterface interface{}, ctx context.Context, decoder func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				request := new(ListAnalysesRequest)
				if err := decoder(request); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return serverInterface.(ExamAnalyzerServiceServer).ListAnalyses(ctx, request)
				}
				serverInfo := &grpc.UnaryServerInfo{
					Server:     serverInterface,
					FullMethod: "/exam_analyzer.v1.ExamAnalyzerService/ListAnalyses",
				}
				handler := func(ctx context.Context, request interface{}) (interface{}, error) {
					return serverInterface.(ExamAnalyzerServiceServer).ListAnalyses(ctx, request.(*ListAnalysesRequest))
				}
				return interceptor(ctx, request, serverInfo, handler)
			},
		},
		{
			MethodName: "GetAnalysis",
			Handler: func(serverInterface interface{}, ctx context.Context, decoder func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				request := new(GetAnalysisRequest)
				if err := decoder(request); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return serverInterface.(ExamAnalyzerServiceServer).GetAnalysis(ctx, request)
				}
				serverInfo := &grpc.UnaryServerInfo{
					Server:     serverInterface,
					FullMethod: "/exam_analyzer.v1.ExamAnalyzerService/GetAnalysis",
				}
				handler := func(ctx context.Context, request interface{}) (interface{}, error) {
					return serverInterface.(ExamAnalyzerServiceServer).GetAnalysis(ctx, request.(*GetAnalysisRequest))
				}
				return interceptor(ctx, request, serverInfo, handler)
			},
		},
		{
			MethodName: "DeleteAnalysis",
			Handler: func(serverInterface interface{}, ctx context.Context, decoder func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				request := new(DeleteAnalysisRequest)
				if err := decoder(request); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return serverInterface.(ExamAnalyzerServiceServer).DeleteAnalysis(ctx, request)
				}
				serverInfo := &grpc.UnaryServerInfo{
					Server:     serverInterface,
					FullMethod: "/exam_analyzer.v1.ExamAnalyzerService/DeleteAnalysis",
				}
				handler := func(ctx context.Context, request interface{}) (interface{}, error) {
					return serverInterface.(ExamAnalyzerServiceServer).DeleteAnalysis(ctx, request.(*DeleteAnalysisRequest))
				}
				return interceptor(ctx, request, serverInfo, handler)
			},
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "exam_analyzer.proto",
}
