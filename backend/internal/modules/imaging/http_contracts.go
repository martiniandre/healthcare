package imaging

import "time"

type HTTPImagingStudyResponse struct {
	ID               string `json:"id"`
	PatientFHIRID    string `json:"patient_fhir_id"`
	Title            string `json:"title"`
	Modality         string `json:"modality"`
	StudyInstanceUID string `json:"study_instance_uid"`
	Status           string `json:"status"`
	DownloadURL      string `json:"download_url,omitempty"`
	CreatedAt        string `json:"created_at"`
}

func NewHTTPImagingStudyResponse(study *ImagingStudy) HTTPImagingStudyResponse {
	return HTTPImagingStudyResponse{
		ID:               study.ID.String(),
		PatientFHIRID:    study.PatientFhirID,
		Title:            study.Title,
		Modality:         study.Modality,
		StudyInstanceUID: study.StudyInstanceUID,
		Status:           study.Status,
		CreatedAt:        study.CreatedAt.Format(time.RFC3339),
	}
}

func NewHTTPImagingStudyResponseWithDownloadURL(study *ImagingStudy, downloadURL string) HTTPImagingStudyResponse {
	studyResponse := NewHTTPImagingStudyResponse(study)
	studyResponse.DownloadURL = downloadURL
	return studyResponse
}

func NewHTTPImagingStudyResponses(studiesList []*ImagingStudy) []HTTPImagingStudyResponse {
	responseList := make([]HTTPImagingStudyResponse, 0, len(studiesList))
	for _, study := range studiesList {
		responseList = append(responseList, NewHTTPImagingStudyResponse(study))
	}
	return responseList
}
