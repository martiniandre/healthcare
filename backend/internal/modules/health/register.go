package health

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func Register(grpcServer *grpc.Server, databaseConnectionPool *pgxpool.Pool, cacheConnectionClient *redis.Client) {
	grpc_health_v1.RegisterHealthServer(grpcServer, NewGRPCHandler(databaseConnectionPool, cacheConnectionClient))
}
