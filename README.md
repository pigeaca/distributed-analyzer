# Distributed Go Analyze Platform
A distributed platform for executing and analyzing tasks — for example, building, statically analyzing, and benchmarking Go projects in a scalable way.

## Architecture Overview

### Components

#### API Gateway
- gRPC/REST entry point
- Authentication (JWT)
- Rate Limiting, Throttling
- Proxy to components

#### Task Service
- Creation and management of tasks
- Statuses (Created, Scheduled, InProgress, Done, Failed)
- Input data validation
- Task storage (PostgresSQL)
- Exposes gRPC API for task operations

#### Scheduler / Dispatcher
- Subscribed to new task queue (Kafka)
- Selects appropriate worker
- Divides task into subtasks if needed
- Assigns to workers
- Exposes gRPC API for scheduling operations

#### Worker Manager
- Manages worker registration
- Monitors worker statuses (healthcheck, capacity)
- Matches worker capabilities (e.g., needs benchmarking support, or race detection)
- Exposes gRPC API for worker management

#### Worker Node (agent)
- Receives tasks via Kafka or gRPC
- Loads model
- Executes task in sandbox (Docker / gVisor)
- Sends result back via gRPC
- Exposes gRPC API for task execution

#### Result Aggregator
- Saves partial results
- Finalizes result when a task is completed
- Sets Done status
- Exposes gRPC API for result operations

#### Billing & Quotas
- Task billing (time, model type, resources)
- Token/balance support
- User quotas
- Exposes gRPC API for billing operations

#### Audit & Logging
- Task tracing, action logging
- Protection against falsification (optional: hash chain, signatures)
- Exposes gRPC API for audit operations

#### Monitoring / Observability
- Prometheus + Grafana
- OpenTelemetry + Jaeger
- Logs via Loki/FluentBit

### Storage
- PostgreSQL: tasks, users, billing
- Redis: states, heartbeats, cache
- MinIO (S3): files (input/output), models

### Communication
- **gRPC**: Synchronous service-to-service communication
- **Kafka**: Asynchronous event-driven communication

## Workflows

### Task Placement Flow
1. User sends a task → API Gateway (REST)
2. Gateway → TaskService (gRPC) → validation, task saving
3. TaskService publishes TaskCreatedEvent → Kafka → Scheduler subscribes

### Task Execution Flow
1. Scheduler selects a worker (by capability) using WorkerManager (gRPC)
2. Scheduler publishes TaskAssignedEvent → Kafka → Worker subscribes
3. Worker executes a task and writes a partial result to S3/MinIO
4. Worker calls ResultAggregator (gRPC) to save partial results
5. Worker publishes SubTaskCompletedEvent → Kafka

### Completion Flow
1. ResultAggregator finalizes a result when all subtasks are completed
2. ResultAggregator publishes TaskCompletedEvent → Kafka
3. TaskService updates task status (subscribes to TaskCompletedEvent)
4. BillingService charges cost (subscribes to TaskCompletedEvent)
5. User can download a result via API Gateway (REST)

### Monitoring Flow
1. All components emit metrics to Prometheus
2. All components send traces to Jaeger via OpenTelemetry
3. All components log to centralized logging (Loki/FluentBit)
4. AuditService records all significant actions (subscribes to AuditEvent)

## Protocol Buffers and gRPC

### Proto Files
The project uses Protocol Buffers (protobuf) for defining service interfaces and message types. Proto files are located in the `proto` directory:

- `proto/task/task.proto`: Task service interface and message types
- `proto/scheduler/scheduler.proto`: Scheduler service interface and message types
- `proto/worker/worker.proto`: Worker service interface and message types
- `proto/result/result.proto`: Result service interface and message types
- `proto/billing/billing.proto`: Billing service interface and message types
- `proto/audit/audit.proto`: Audit service interface and message types
- `proto/user/user.proto`: User service interface and message types
- `proto/kafka/messages.proto`: Kafka message types for asynchronous communication

### Code Generation
To generate Go code from proto files, use the following command:

```bash
./scripts/generate.sh
```

This will generate:
- gRPC service interfaces and clients
- Message type structs
- Serialization/deserialization code

The generated code will be placed in the `pkg/proto` directory.
