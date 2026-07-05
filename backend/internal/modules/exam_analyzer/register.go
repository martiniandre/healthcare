package exam_analyzer

import (
	"github.com/jackc/pgx/v5/pgxpool"
	examanalyzerpb "github.com/healthcare/backend/internal/modules/exam_analyzer/pb"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"google.golang.org/grpc"
)

type Dependency struct {
	DB          *pgxpool.Pool
	ProjectID   string
	LocationID  string
	VertexModel string
	EventBus    eventbus.Bus
}

func Register(grpcServer *grpc.Server, dep Dependency) (Repository, Service, *Worker) {
	repo := NewRepository(dep.DB)
	svc := NewService(repo, dep.ProjectID, dep.LocationID, dep.VertexModel)
	handler := NewGRPCHandler(repo, svc)
	examanalyzerpb.RegisterExamAnalyzerServiceServer(grpcServer, handler)
	worker := NewWorker(repo, svc, dep.EventBus)
	return repo, svc, worker
}
