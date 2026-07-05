package pb

import "context"

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

func RegisterExamAnalyzerServiceServer(_ interface{}, server ExamAnalyzerServiceServer) {}
