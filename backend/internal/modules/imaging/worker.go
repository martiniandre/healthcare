package imaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/redis/go-redis/v9"
)

const imagingStudyResourceType = "ImagingStudy"

type Worker struct {
	dbRepository     Repository
	redisClient      *redis.Client
	healthcareClient healthcare.FHIRClient
	stopChannel      chan struct{}
	stopOnce         sync.Once
}

func NewWorker(dbRepository Repository, redisClient *redis.Client, healthcareClient healthcare.FHIRClient) *Worker {
	return &Worker{
		dbRepository:     dbRepository,
		redisClient:      redisClient,
		healthcareClient: healthcareClient,
		stopChannel:      make(chan struct{}),
	}
}

func (worker *Worker) Start(ctx context.Context) {
	slog.Info("imaging background worker started successfully")
	for {
		select {
		case <-worker.stopChannel:
			slog.Info("imaging background worker received shutdown signal")
			return
		case <-ctx.Done():
			slog.Info("imaging background worker context cancelled")
			return
		default:
			if worker.redisClient == nil {
				time.Sleep(5 * time.Second)
				continue
			}

			popResults, popError := worker.redisClient.BRPop(ctx, 5*time.Second, processingQueueKey).Result()
			if popError != nil {
				if !errors.Is(popError, redis.Nil) && !errors.Is(popError, context.Canceled) {
					slog.Warn("failed to pop imaging job from queue", "error", popError)
				}
				continue
			}

			if len(popResults) < 2 {
				continue
			}

			studyIDString := popResults[1]
			worker.processDICOM(ctx, studyIDString)
		}
	}
}

func (worker *Worker) Stop() {
	worker.stopOnce.Do(func() {
		close(worker.stopChannel)
	})
}

func (worker *Worker) processDICOM(ctx context.Context, studyIDString string) {
	studyID, parseError := uuid.Parse(studyIDString)
	if parseError != nil {
		slog.Error("failed to parse study id from queue", "id", studyIDString, "error", parseError)
		return
	}

	study, dbError := worker.dbRepository.GetImagingStudy(ctx, studyID)
	if dbError != nil {
		slog.Error("failed to fetch operational study from db", "id", studyIDString, "error", dbError)
		return
	}

	study.Status = ImagingStudyStatusProcessing
	study.UpdatedAt = time.Now()
	if updateError := worker.dbRepository.UpdateImagingStudy(ctx, study); updateError != nil {
		slog.Error("failed to mark imaging study as processing", "id", studyIDString, "error", updateError)
		return
	}

	imagingSeries, seriesError := worker.buildImagingSeries(ctx, study)
	if seriesError != nil {
		slog.Error("failed to load imaging study instances", "id", studyIDString, "error", seriesError)
		study.Status = ImagingStudyStatusFailed
		study.UpdatedAt = time.Now()
		if updateError := worker.dbRepository.UpdateImagingStudy(ctx, study); updateError != nil {
			slog.Error("failed to persist failed imaging study status", "id", studyIDString, "error", updateError)
		}
		return
	}

	startedAt := study.CreatedAt
	if study.StudyDateTime != nil {
		startedAt = *study.StudyDateTime
	}

	fhirStudy := fhir.NewImagingStudyResource(
		study.PatientFhirID,
		"available",
		startedAt.Format(time.RFC3339),
		study.Title,
		imagingSeries,
	)

	var registrationError error
	if worker.healthcareClient != nil {
		registrationError = worker.registerImagingStudy(ctx, study, fhirStudy)
	}

	if registrationError != nil {
		slog.Error("failed to register ImagingStudy FHIR resource", "id", studyIDString, "error", registrationError)
		study.Status = ImagingStudyStatusFailed
	} else {
		study.Status = ImagingStudyStatusReady
	}

	study.UpdatedAt = time.Now()
	if updateError := worker.dbRepository.UpdateImagingStudy(ctx, study); updateError != nil {
		slog.Error("failed to persist processed imaging study status", "id", studyIDString, "status", study.Status, "error", updateError)
		return
	}
	slog.Info("imaging study processed and registered successfully", "id", studyIDString, "status", study.Status)
}

func (worker *Worker) registerImagingStudy(ctx context.Context, study *ImagingStudy, fhirStudy *fhir.ImagingStudyResource) error {
	if study.FHIRResourceID != "" {
		fhirStudy.ID = study.FHIRResourceID
		_, updateError := worker.healthcareClient.UpdateResource(ctx, imagingStudyResourceType, study.FHIRResourceID, fhirStudy)
		return updateError
	}

	responseBody, createError := worker.healthcareClient.CreateResource(ctx, imagingStudyResourceType, fhirStudy)
	if createError != nil {
		return createError
	}
	var createdResource struct {
		ID string `json:"id"`
	}
	if decodeError := json.Unmarshal(responseBody, &createdResource); decodeError != nil {
		return fmt.Errorf("failed to read created ImagingStudy id: %w", decodeError)
	}
	study.FHIRResourceID = createdResource.ID
	return nil
}

func (worker *Worker) buildImagingSeries(ctx context.Context, study *ImagingStudy) ([]fhir.Series, error) {
	instances, listError := worker.dbRepository.ListInstances(ctx, study.ID)
	if listError != nil {
		return nil, listError
	}

	var seriesOrder []string
	seriesByUID := make(map[string]*fhir.Series)
	for _, instance := range instances {
		series, known := seriesByUID[instance.SeriesInstanceUID]
		if !known {
			series = &fhir.Series{
				Uid:    instance.SeriesInstanceUID,
				Number: instance.SeriesNumber,
				Modality: fhir.Coding{
					System:  "http://dicom.nema.org/resources/ontology/DCM",
					Code:    instance.Modality,
					Display: instance.Modality,
				},
			}
			seriesByUID[instance.SeriesInstanceUID] = series
			seriesOrder = append(seriesOrder, instance.SeriesInstanceUID)
		}
		series.Instance = append(series.Instance, fhir.Instance{
			Uid: instance.SOPInstanceUID,
			SopClass: fhir.Coding{
				System: "urn:ietf:rfc:3986",
				Code:   "urn:oid:" + instance.SOPClassUID,
			},
			Number: instance.InstanceNumber,
		})
	}

	imagingSeries := make([]fhir.Series, 0, len(seriesOrder))
	for _, seriesUID := range seriesOrder {
		imagingSeries = append(imagingSeries, *seriesByUID[seriesUID])
	}
	return imagingSeries, nil
}
