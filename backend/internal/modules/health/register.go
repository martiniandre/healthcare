package health

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Dependency struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func Register(grpcServer *grpc.Server, dep Dependency) {
	grpc_health_v1.RegisterHealthServer(grpcServer, NewGRPCHandler(dep.DB, dep.Redis))
}
