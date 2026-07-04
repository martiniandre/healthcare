package encounter

import (
	encounterpb "github.com/healthcare/backend/internal/modules/encounter/pb"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"google.golang.org/grpc"
)

type Dependency struct {
	FHIRClient healthcare.FHIRClient
	EventBus   eventbus.Bus
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.FHIRClient)
	svc := NewService(repo, dep.EventBus)
	handler := NewGRPCHandler(svc)
	encounterpb.RegisterEncounterServiceServer(grpcServer, handler)
	return svc
}
