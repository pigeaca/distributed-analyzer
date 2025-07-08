package bootstrap

import (
	"fmt"
	"github.com/distributedmarketplace/internal/task/config"
	"github.com/distributedmarketplace/internal/task/grpc"
	"github.com/distributedmarketplace/internal/task/kafka/handler"
	tProducer "github.com/distributedmarketplace/internal/task/kafka/producer"
	"github.com/distributedmarketplace/internal/task/service"
	"github.com/distributedmarketplace/pkg/application"
	app3 "github.com/distributedmarketplace/pkg/application/grpc"
	app "github.com/distributedmarketplace/pkg/application/kafka"
	"github.com/distributedmarketplace/pkg/kafka"
	pb "github.com/distributedmarketplace/pkg/proto/task"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func StartApplication(cfg config.Config) error {
	taskService := service.NewTaskServiceImpl()
	producer, kafkaComponent := initKafka(cfg, taskService)
	_, grpcComponent := initGrpc(cfg, producer, taskService)
	runner := application.NewApplicationRunner(grpcComponent, kafkaComponent)
	return runner.StartBlocking()
}

func initGrpc(cfg config.Config, producer *kafka.Producer, service service.TaskService) (*stdgrpc.Server, *app3.GrpcComponent) {
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

func initKafka(cfg config.Config, taskService service.TaskService) (*kafka.Producer, *app.KafkaComponent) {
	producer := kafka.NewProducer(cfg.KafkaBrokers)
	taskHandler := handler.NewTaskMessageHandler(taskService)
	topics := []string{"task-status-changed", "task-completed", "task-failed"}
	consumer := kafka.NewConsumer(topics, cfg.KafkaBrokers, cfg.KafkaGroupID, taskHandler)
	return producer, app.NewKafkaComponent(consumer)
}
