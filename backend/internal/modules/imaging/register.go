package imaging

import (
	"github.com/healthcare/backend/internal/modules/imaging/pb"
	"github.com/healthcare/backend/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Dependency struct {
	DB         *pgxpool.Pool
	Storage    storage.StorageClient
	Redis      *redis.Client
	BucketName string
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.DB)
	svc := NewService(repo, dep.Storage, dep.Redis, dep.BucketName)
	handler := NewGRPCHandler(svc)
	pb.RegisterImagingServiceServer(grpcServer, handler)
	return svc
}
