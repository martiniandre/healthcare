package imaging

import (
	"context"
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

	studyInstanceUID := fmt.Sprintf("1.2.826.0.1.3680043.8.498.%s", uuid.New().String())

	imagingSeries := []fhir.Series{
		{
			Uid:    fmt.Sprintf("1.2.826.0.1.3680043.8.498.series.%s", uuid.New().String()),
			Number: 1,
			Modality: fhir.Coding{
				System:  "http://dicom.nema.org/resources/ontology/DCM",
				Code:    study.Modality,
				Display: study.Modality,
			},
			Instance: []fhir.Instance{
				{
					Uid: studyInstanceUID,
					SopClass: fhir.Coding{
						System: "http://dicom.nema.org/resources/ontology/DCM",
						Code:   "1.2.840.10008.5.1.4.1.1.7",
					},
					Number: 1,
				},
			},
		},
	}

	fhirStudy := fhir.NewImagingStudyResource(
		study.PatientFhirID,
		"available",
		study.CreatedAt.Format(time.RFC3339),
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
		study.Status = ImagingStudyStatusProcessed
		study.StudyInstanceUID = studyInstanceUID
	}

	study.UpdatedAt = time.Now()
	if updateError := worker.dbRepository.UpdateImagingStudy(ctx, study); updateError != nil {
		slog.Error("failed to persist processed imaging study status", "id", studyIDString, "status", study.Status, "error", updateError)
		return
	}
	slog.Info("imaging study processed and registered successfully", "id", studyIDString, "status", study.Status)
}
