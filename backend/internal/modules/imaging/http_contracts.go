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

type HTTPUploadWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HTTPRejectedInstance struct {
	FileName string `json:"fileName"`
	Code     string `json:"code"`
	Reason   string `json:"reason"`
}

type HTTPUploadStudyResponse struct {
	ID                string                 `json:"id"`
	PatientFHIRID     string                 `json:"patient_fhir_id"`
	Title             string                 `json:"title"`
	Modality          string                 `json:"modality"`
	StudyInstanceUID  string                 `json:"study_instance_uid"`
	Status            string                 `json:"status"`
	CreatedAt         string                 `json:"created_at"`
	StudyInstanceUid  string                 `json:"studyInstanceUid"`
	StudyDate         string                 `json:"studyDate,omitempty"`
	AccessionNumber   string                 `json:"accessionNumber,omitempty"`
	SeriesInstanceUid string                 `json:"seriesInstanceUid"`
	SopInstanceUid    string                 `json:"sopInstanceUid"`
	Warnings          []HTTPUploadWarning    `json:"warnings"`
	StudyID           string                 `json:"studyId"`
	StudyCreated      bool                   `json:"studyCreated"`
	InstancesAccepted int                    `json:"instancesAccepted"`
	InstancesRejected []HTTPRejectedInstance `json:"instancesRejected"`
}

func NewHTTPUploadStudyResponse(uploadResult *UploadDICOMResult) HTTPUploadStudyResponse {
	study := uploadResult.Study
	warnings := make([]HTTPUploadWarning, 0, len(uploadResult.Warnings))
	for _, warning := range uploadResult.Warnings {
		warnings = append(warnings, HTTPUploadWarning{Code: warning.Code, Message: warning.MessageKey})
	}

	rejections := make([]HTTPRejectedInstance, 0, len(uploadResult.InstancesRejected))
	for _, rejection := range uploadResult.InstancesRejected {
		rejections = append(rejections, HTTPRejectedInstance{FileName: rejection.FileName, Code: rejection.Code, Reason: rejection.Reason})
	}

	studyDate := ""
	if study.StudyDateTime != nil {
		studyDate = study.StudyDateTime.Format(time.RFC3339)
	}

	return HTTPUploadStudyResponse{
		ID:                study.ID.String(),
		PatientFHIRID:     study.PatientFhirID,
		Title:             study.Title,
		Modality:          study.Modality,
		StudyInstanceUID:  study.StudyInstanceUID,
		Status:            study.Status,
		CreatedAt:         study.CreatedAt.Format(time.RFC3339),
		StudyInstanceUid:  study.StudyInstanceUID,
		StudyDate:         studyDate,
		AccessionNumber:   study.AccessionNumber,
		SeriesInstanceUid: uploadResult.SeriesInstanceUID,
		SopInstanceUid:    uploadResult.SOPInstanceUID,
		Warnings:          warnings,
		StudyID:           study.ID.String(),
		StudyCreated:      uploadResult.StudyCreated,
		InstancesAccepted: uploadResult.InstancesAccepted,
		InstancesRejected: rejections,
	}
}

type HTTPStudyDownloadURLResponse struct {
	DownloadURL string `json:"downloadUrl"`
	ExpiresAt   string `json:"expiresAt"`
}

type HTTPStudyInstanceResponse struct {
	ID                string `json:"id"`
	SOPInstanceUID    string `json:"sopInstanceUid"`
	SeriesInstanceUID string `json:"seriesInstanceUid"`
	InstanceNumber    *int   `json:"instanceNumber"`
	FileName          string `json:"fileName"`
	DownloadURL       string `json:"downloadUrl"`
	ExpiresAt         string `json:"expiresAt"`
}

func NewHTTPStudyInstanceResponse(download StudyInstanceDownload) HTTPStudyInstanceResponse {
	var instanceNumber *int
	if download.Instance.InstanceNumber > 0 {
		number := download.Instance.InstanceNumber
		instanceNumber = &number
	}
	return HTTPStudyInstanceResponse{
		ID:                download.Instance.ID.String(),
		SOPInstanceUID:    download.Instance.SOPInstanceUID,
		SeriesInstanceUID: download.Instance.SeriesInstanceUID,
		InstanceNumber:    instanceNumber,
		FileName:          download.Instance.FileName,
		DownloadURL:       download.DownloadURL,
		ExpiresAt:         download.ExpiresAt.Format(time.RFC3339),
	}
}

func NewHTTPStudyInstanceResponses(downloads []StudyInstanceDownload) []HTTPStudyInstanceResponse {
	responseList := make([]HTTPStudyInstanceResponse, 0, len(downloads))
	for _, download := range downloads {
		responseList = append(responseList, NewHTTPStudyInstanceResponse(download))
	}
	return responseList
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
