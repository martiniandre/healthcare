package imaging

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/imaging/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) UploadDICOM(stream pb.ImagingService_UploadDICOMServer) error {
	firstRequestMsg, receiveError := stream.Recv()
	if receiveError != nil {
		return mapImagingError(receiveError)
	}

	metadataMsg := firstRequestMsg.GetMetadata()
	if metadataMsg == nil {
		return apperrors.ErrBadRequest.ToGRPC()
	}

	pipeReader, pipeWriter := io.Pipe()
	var totalBytesUploaded int64

	go func() {
		defer pipeWriter.Close()
		for {
			chunkRequestMsg, streamError := stream.Recv()
			if errors.Is(streamError, io.EOF) {
				return
			}
			if streamError != nil {
				_ = pipeWriter.CloseWithError(streamError)
				return
			}

			chunkBytes := chunkRequestMsg.GetChunk()
			if chunkBytes != nil {
				_, writeError := pipeWriter.Write(chunkBytes)
				if writeError != nil {
					_ = pipeWriter.CloseWithError(writeError)
					return
				}

				totalBytesUploaded += int64(len(chunkBytes))
				_ = stream.Send(&pb.UploadDICOMResponse{
					ImagingStudyId: "",
					Status:         "UPLOADING",
					BytesUploaded:  totalBytesUploaded,
				})
			}
		}
	}()

	study, uploadError := handler.service.UploadDICOMStream(
		stream.Context(),
		metadataMsg.PatientFhirId,
		metadataMsg.Title,
		metadataMsg.Modality,
		pipeReader,
	)
	if uploadError != nil {
		return mapImagingError(uploadError)
	}

	response := &pb.UploadDICOMResponse{
		ImagingStudyId: study.ID.String(),
		Status:         study.Status,
		BytesUploaded:  totalBytesUploaded,
	}

	return stream.Send(response)
}

func (handler *GRPCHandler) GetImagingStudy(ctx context.Context, req *pb.GetImagingStudyRequest) (*pb.GetImagingStudyResponse, error) {
	studyID, parseError := uuid.Parse(req.ImagingStudyId)
	if parseError != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	study, queryError := handler.service.GetImagingStudy(ctx, studyID)
	if queryError != nil {
		return nil, mapImagingError(queryError)
	}

	return &pb.GetImagingStudyResponse{
		ImagingStudyId:   study.ID.String(),
		PatientFhirId:    study.PatientFhirID,
		Title:            study.Title,
		Modality:         study.Modality,
		StudyInstanceUid: study.StudyInstanceUID,
		Status:           study.Status,
		CreatedAt:        study.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (handler *GRPCHandler) ListImagingStudies(ctx context.Context, req *pb.ListImagingStudiesRequest) (*pb.ListImagingStudiesResponse, error) {
	studies, queryError := handler.service.ListImagingStudies(ctx, req.PatientFhirId)
	if queryError != nil {
		return nil, mapImagingError(queryError)
	}

	studiesResponses := make([]*pb.GetImagingStudyResponse, 0, len(studies))
	for _, study := range studies {
		studiesResponses = append(studiesResponses, &pb.GetImagingStudyResponse{
			ImagingStudyId:   study.ID.String(),
			PatientFhirId:    study.PatientFhirID,
			Title:            study.Title,
			Modality:         study.Modality,
			StudyInstanceUid: study.StudyInstanceUID,
			Status:           study.Status,
			CreatedAt:        study.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListImagingStudiesResponse{Studies: studiesResponses}, nil
}

func (handler *GRPCHandler) GetDICOMDownloadURL(ctx context.Context, req *pb.GetDICOMDownloadURLRequest) (*pb.GetDICOMDownloadURLResponse, error) {
	studyID, parseError := uuid.Parse(req.ImagingStudyId)
	if parseError != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	downloadURL, expiresAt, queryError := handler.service.GetDownloadURL(ctx, studyID)
	if queryError != nil {
		return nil, mapImagingError(queryError)
	}

	return &pb.GetDICOMDownloadURLResponse{
		DownloadUrl: downloadURL,
		ExpiresAt:   expiresAt.Format(time.RFC3339),
	}, nil
}

func mapImagingError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrImagingStudyNotFound):
		return apperrors.ErrImagingStudyNotFound.ToGRPC()
	case errors.Is(err, ErrInvalidDICOM):
		return apperrors.ErrInvalidDICOM.ToGRPC()
	case errors.Is(err, ErrDICOMTooLarge):
		return apperrors.ErrRateLimitExceeded.ToGRPC()
	default:
		return apperrors.ErrInternalServer.ToGRPC()
	}
}
