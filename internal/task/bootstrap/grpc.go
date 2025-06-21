package bootstrap

import (
	"fmt"
	"github.com/distributedmarketplace/internal/task/config"
	"github.com/distributedmarketplace/internal/task/grpc"
	tProducer "github.com/distributedmarketplace/internal/task/kafka/producer"
	"github.com/distributedmarketplace/internal/task/service"
	app3 "github.com/distributedmarketplace/pkg/application/grpc"
	"github.com/distributedmarketplace/pkg/kafka"
	pb "github.com/distributedmarketplace/pkg/proto/task"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func InitGrpc(cfg config.Config, producer *kafka.Producer, service service.TaskService) (*stdgrpc.Server, *app3.GrpcComponent) {
	grpcServer := stdgrpc.NewServer()
	taskProducer := tProducer.NewTaskProducer(producer)
	taskGrpcServer := grpc.NewTaskServer(service, taskProducer)
	pb.RegisterTaskServiceServer(grpcServer, taskGrpcServer)
	reflection.Register(grpcServer)

	grpcAddr := fmt.Sprintf(":%s", cfg.GrpcPort)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}
	return grpcServer, app3.NewGrpcComponent(grpcServer, listener)
}
