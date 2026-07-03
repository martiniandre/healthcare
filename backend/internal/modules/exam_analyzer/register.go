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

var (
	WorkerInstance *Worker
	Repo           Repository
	Svc            Service
)

func Register(dep Dependency) {
	Repo = NewRepository(dep.DB)
	Svc = NewService(Repo, dep.ProjectID, dep.LocationID, dep.VertexModel)
	WorkerInstance = NewWorker(Repo, Svc)
}
