package main

import (
	"net"

	"github.com/najimovmashhurbek/project-api/user-service.ozim/config"
	pb "github.com/najimovmashhurbek/project-api/user-service.ozim/genproto"
	"github.com/najimovmashhurbek/project-api/user-service.ozim/pkg/db"
	"github.com/najimovmashhurbek/project-api/user-service.ozim/pkg/logger"
	"github.com/najimovmashhurbek/project-api/user-service.ozim/service"
	grpcClient "github.com/najimovmashhurbek/project-api/user-service.ozim/service/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "user.ozim")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}
	grpc1, _ := grpcClient.New(cfg)
	userService := service.NewUserService(connDB, log, grpc1)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
