package imaging

import (
	"github.com/healthcare/backend/internal/modules/imaging/pb"
	"github.com/healthcare/backend/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func Register(grpcServer *grpc.Server, dbPool *pgxpool.Pool, storageClient storage.StorageClient, redisClient *redis.Client, bucketName string) Service {
	repo := NewRepository(dbPool)
	svc := NewService(repo, storageClient, redisClient, bucketName)
	handler := NewGRPCHandler(svc)
	pb.RegisterImagingServiceServer(grpcServer, handler)
	return svc
}
