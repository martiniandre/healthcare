package pb

import (
	"context"

	"google.golang.org/grpc"
)

type ImagingServiceServer interface {
	UploadDICOM(stream ImagingService_UploadDICOMServer) error
	GetImagingStudy(ctx context.Context, req *GetImagingStudyRequest) (*GetImagingStudyResponse, error)
	ListImagingStudies(ctx context.Context, req *ListImagingStudiesRequest) (*ListImagingStudiesResponse, error)
	GetDICOMDownloadURL(ctx context.Context, req *GetDICOMDownloadURLRequest) (*GetDICOMDownloadURLResponse, error)
}

type ImagingService_UploadDICOMServer interface {
	Send(response *UploadDICOMResponse) error
	Recv() (*UploadDICOMRequest, error)
	grpc.ServerStream
}

type UploadDICOMRequest struct {
	Payload isUploadDICOMRequest_Payload
}

type isUploadDICOMRequest_Payload interface {
	isUploadDICOMRequest_Payload()
}

type UploadDICOMRequest_Metadata struct {
	Metadata *UploadMetadata
}

type UploadDICOMRequest_Chunk struct {
	Chunk []byte
}

func (payloadMetadata *UploadMetadata) GetPatientFhirId() string {
	return payloadMetadata.PatientFhirId
}

func (payloadMetadata *UploadMetadata) GetTitle() string {
	return payloadMetadata.Title
}

func (payloadMetadata *UploadMetadata) GetModality() string {
	return payloadMetadata.Modality
}

func (uploadRequest *UploadDICOMRequest) GetMetadata() *UploadMetadata {
	if payloadMetadata, ok := uploadRequest.Payload.(*UploadDICOMRequest_Metadata); ok {
		return payloadMetadata.Metadata
	}
	return nil
}

func (uploadRequest *UploadDICOMRequest) GetChunk() []byte {
	if payloadChunk, ok := uploadRequest.Payload.(*UploadDICOMRequest_Chunk); ok {
		return payloadChunk.Chunk
	}
	return nil
}

func (requestMetadata *UploadDICOMRequest_Metadata) isUploadDICOMRequest_Payload() {}
func (requestChunk *UploadDICOMRequest_Chunk) isUploadDICOMRequest_Payload()       {}

type UploadMetadata struct {
	PatientFhirId string
	Title         string
	Modality      string
}

type UploadDICOMResponse struct {
	ImagingStudyId string
	Status         string
	BytesUploaded  int64
}

type GetImagingStudyRequest struct {
	ImagingStudyId string
}

type GetImagingStudyResponse struct {
	ImagingStudyId   string
	PatientFhirId    string
	Title            string
	Modality         string
	StudyInstanceUid string
	Status           string
	CreatedAt        string
}

type ListImagingStudiesRequest struct {
	PatientFhirId string
}

type ListImagingStudiesResponse struct {
	Studies []*GetImagingStudyResponse
}

type GetDICOMDownloadURLRequest struct {
	ImagingStudyId string
}

type GetDICOMDownloadURLResponse struct {
	DownloadUrl string
	ExpiresAt   string
}

func RegisterImagingServiceServer(server *grpc.Server, handler ImagingServiceServer) {
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "imaging.v1.ImagingService",
		HandlerType: (*ImagingServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetImagingStudy",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(GetImagingStudyRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return srv.(ImagingServiceServer).GetImagingStudy(ctx, in)
					}
					info := &grpc.UnaryServerInfo{
						Server:     srv,
						FullMethod: "/imaging.v1.ImagingService/GetImagingStudy",
					}
					handlerFunc := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(ImagingServiceServer).GetImagingStudy(ctx, req.(*GetImagingStudyRequest))
					}
					return interceptor(ctx, in, info, handlerFunc)
				},
			},
			{
				MethodName: "ListImagingStudies",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(ListImagingStudiesRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return srv.(ImagingServiceServer).ListImagingStudies(ctx, in)
					}
					info := &grpc.UnaryServerInfo{
						Server:     srv,
						FullMethod: "/imaging.v1.ImagingService/ListImagingStudies",
					}
					handlerFunc := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(ImagingServiceServer).ListImagingStudies(ctx, req.(*ListImagingStudiesRequest))
					}
					return interceptor(ctx, in, info, handlerFunc)
				},
			},
			{
				MethodName: "GetDICOMDownloadURL",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(GetDICOMDownloadURLRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return srv.(ImagingServiceServer).GetDICOMDownloadURL(ctx, in)
					}
					info := &grpc.UnaryServerInfo{
						Server:     srv,
						FullMethod: "/imaging.v1.ImagingService/GetDICOMDownloadURL",
					}
					handlerFunc := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(ImagingServiceServer).GetDICOMDownloadURL(ctx, req.(*GetDICOMDownloadURLRequest))
					}
					return interceptor(ctx, in, info, handlerFunc)
				},
			},
		},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "UploadDICOM",
				Handler:       func(srv interface{}, stream grpc.ServerStream) error {
					return srv.(ImagingServiceServer).UploadDICOM(&imagingServiceUploadDICOMServer{stream})
				},
				ServerStreams: true,
				ClientStreams: true,
			},
		},
		Metadata: "proto/imaging.proto",
	}, handler)
}

type imagingServiceUploadDICOMServer struct {
	grpc.ServerStream
}

func (serverStream *imagingServiceUploadDICOMServer) Send(response *UploadDICOMResponse) error {
	return serverStream.SendMsg(response)
}

func (serverStream *imagingServiceUploadDICOMServer) Recv() (*UploadDICOMRequest, error) {
	requestPayload := new(UploadDICOMRequest)
	if receiveError := serverStream.RecvMsg(requestPayload); receiveError != nil {
		return nil, receiveError
	}
	return requestPayload, nil
}
