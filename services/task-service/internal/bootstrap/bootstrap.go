package bootstrap

import (
	"fmt"
	"github.com/pigeaca/DistributedMarketplace/libs/application"
	grpcApp "github.com/pigeaca/DistributedMarketplace/libs/application/grpc"
	kafkaApp "github.com/pigeaca/DistributedMarketplace/libs/application/kafka"
	"github.com/pigeaca/DistributedMarketplace/libs/kafka"
	pb "github.com/pigeaca/DistributedMarketplace/libs/proto/task"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/config"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/grpc"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/kafka/handler"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/kafka/producer"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/service"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func StartApplication(cfg *config.Config) error {
	taskService := service.NewTaskServiceImpl()
	kafkaProducer, kafkaComponent := initKafka(cfg, taskService)
	_, grpcComponent := initGrpc(cfg, kafkaProducer, taskService)
	runner := application.NewApplicationRunner(grpcComponent, kafkaComponent)
	return runner.Start()
}

func initGrpc(cfg *config.Config, kafkaProducer *kafka.Producer, service service.TaskService) (*stdgrpc.Server, *grpcApp.GrpcComponent) {
	grpcServer := stdgrpc.NewServer()
	taskProducer := producer.NewTaskProducer(kafkaProducer)
	taskGrpcServer := grpc.NewTaskServer(service, taskProducer)
	pb.RegisterTaskServiceServer(grpcServer, taskGrpcServer)
	reflection.Register(grpcServer)

	grpcAddr := fmt.Sprintf(":%s", cfg.ServerConfig)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}
	return grpcServer, grpcApp.NewGrpcComponent(grpcServer, listener)
}

func initKafka(cfg *config.Config, taskService service.TaskService) (*kafka.Producer, *kafkaApp.KafkaComponent) {
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers)
	taskHandler := handler.NewTaskMessageHandler(taskService)
	topics := []string{"task-status-changed", "task-completed", "task-failed"}
	consumer := kafka.NewConsumer(topics, cfg.Kafka.Brokers, cfg.Kafka.GroupID, taskHandler)
	return kafkaProducer, kafkaApp.NewKafkaComponent(consumer)
}
