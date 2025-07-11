package bootstrap

import (
	"fmt"
	"github.com/pigeaca/DistributedMarketplace/libs/application"
	app2 "github.com/pigeaca/DistributedMarketplace/libs/application/grpc"
	app "github.com/pigeaca/DistributedMarketplace/libs/application/kafka"
	"github.com/pigeaca/DistributedMarketplace/libs/kafka"
	pb "github.com/pigeaca/DistributedMarketplace/libs/proto/task"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/config"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/grpc"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/kafka/handler"
	producer2 "github.com/pigeaca/DistributedMarketplace/services/task-service/internal/kafka/producer"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/service"
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

func initGrpc(cfg config.Config, producer *kafka.Producer, service service.TaskService) (*stdgrpc.Server, *app2.GrpcComponent) {
	grpcServer := stdgrpc.NewServer()
	taskProducer := producer2.NewTaskProducer(producer)
	taskGrpcServer := grpc.NewTaskServer(service, taskProducer)
	pb.RegisterTaskServiceServer(grpcServer, taskGrpcServer)
	reflection.Register(grpcServer)

	grpcAddr := fmt.Sprintf(":%s", cfg.ServerConfig)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}
	return grpcServer, app2.NewGrpcComponent(grpcServer, listener)
}

func initKafka(cfg config.Config, taskService service.TaskService) (*kafka.Producer, *app.KafkaComponent) {
	producer := kafka.NewProducer(cfg.Kafka.Brokers)
	taskHandler := handler.NewTaskMessageHandler(taskService)
	topics := []string{"task-status-changed", "task-completed", "task-failed"}
	consumer := kafka.NewConsumer(topics, cfg.Kafka.Brokers, cfg.Kafka.GroupID, taskHandler)
	return producer, app.NewKafkaComponent(consumer)
}
