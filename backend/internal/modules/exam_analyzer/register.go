package exam_analyzer

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(databasePool *pgxpool.Pool, projectID, locationID, vertexModel string) (Repository, Service, *Worker) {
	repositoryInstance := NewRepository(databasePool)
	serviceInstance := NewService(projectID, locationID, vertexModel)
	workerInstance := NewWorker(repositoryInstance, serviceInstance)
	return repositoryInstance, serviceInstance, workerInstance
}
