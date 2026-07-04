package exam_analyzer

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependency struct {
	DB          *pgxpool.Pool
	ProjectID   string
	LocationID  string
	VertexModel string
}

func Register(dep Dependency) (Repository, Service, *Worker) {
	repo := NewRepository(dep.DB)
	svc := NewService(repo, dep.ProjectID, dep.LocationID, dep.VertexModel)
	worker := NewWorker(repo, svc)
	return repo, svc, worker
}
