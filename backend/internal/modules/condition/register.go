package condition

import (
	conditionpb "github.com/healthcare/backend/internal/modules/condition/pb"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"google.golang.org/grpc"
)

type Dependency struct {
	FHIRClient healthcare.FHIRClient
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.FHIRClient)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	conditionpb.RegisterConditionServiceServer(grpcServer, handler)
	return svc
}
