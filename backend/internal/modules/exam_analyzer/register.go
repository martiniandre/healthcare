package exam_analyzer

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(databasePool *pgxpool.Pool, projectID, locationID string) (Repository, Service, *Worker) {
	repositoryInstance := NewRepository(databasePool)
	serviceInstance := NewService(projectID, locationID)
	workerInstance := NewWorker(repositoryInstance, serviceInstance)
	return repositoryInstance, serviceInstance, workerInstance
}
