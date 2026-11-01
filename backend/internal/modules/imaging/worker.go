package imaging

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/redis/go-redis/v9"
)

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

			popResults, popError := worker.redisClient.BRPop(ctx, 5*time.Second, "dicom_processing_queue").Result()
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

	imagingSeries := worker.buildImagingSeries(ctx, study)

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

	var creationError error
	if worker.healthcareClient != nil {
		_, creationError = worker.healthcareClient.CreateResource(ctx, "ImagingStudy", fhirStudy)
	}

	if creationError != nil {
		slog.Error("failed to register ImagingStudy FHIR resource", "id", studyIDString, "error", creationError)
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

func (worker *Worker) buildImagingSeries(ctx context.Context, study *ImagingStudy) []fhir.Series {
	instanceIndex, indexError := worker.dbRepository.GetInstanceIndex(ctx, study.ID)
	if indexError != nil {
		slog.Warn("imaging study has no instance index", "id", study.ID.String(), "error", indexError)
		return nil
	}

	return []fhir.Series{
		{
			Uid:    instanceIndex.SeriesInstanceUID,
			Number: instanceIndex.SeriesNumber,
			Modality: fhir.Coding{
				System:  "http://dicom.nema.org/resources/ontology/DCM",
				Code:    study.Modality,
				Display: study.Modality,
			},
			Instance: []fhir.Instance{
				{
					Uid: instanceIndex.SOPInstanceUID,
					SopClass: fhir.Coding{
						System: "urn:ietf:rfc:3986",
						Code:   "urn:oid:" + instanceIndex.SOPClassUID,
					},
					Number: instanceIndex.InstanceNumber,
				},
			},
		},
	}
}
