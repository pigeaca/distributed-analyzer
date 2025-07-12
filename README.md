# Distributed Go Analyze Platform

**Distributed Go Analyze Platform** — a scalable, modular, and extensible system for building, statically analyzing, and
benchmarking Go projects across a distributed infrastructure.

Running CI at scale, validating open-source libraries, or building a Go-aware cloud platform.

---

## Key Features

* High-throughput task execution
* Intelligent worker selection by capabilities (e.g., race detection, benchmark mode)
* Real-time observability and metrics
* gRPC for internal APIs, Kafka for async workflows, REST for external clients
* Designed with modularity, auditability, and cost-tracking in mind

---

## Architecture Overview

### Core Components

#### API Gateway

* Entry point for all clients (REST and gRPC)
* Handles JWT authentication, rate limiting, and request routing

#### Task Service

* Creates and validates tasks
* Stores metadata in PostgreSQL
* Publishes events (e.g. `TaskCreated`) for downstream processing

#### Scheduler / Dispatcher

* Listens to Kafka for new task events
* Selects the optimal worker node using Worker Manager
* Splits tasks into subtasks and assigns them

#### Worker Manager

* Keeps track of available worker nodes
* Matches workers based on capabilities and availability
* Monitors health and load via heartbeats

#### Worker Node (agent)

* Runs inside a sandbox (Docker/gVisor)
* Pulls tasks from Kafka or receives via gRPC
* Executes Go code and reports results

#### Result Aggregator

* Collects partial results from workers
* Finalizes output when all subtasks complete
* Notifies Task Service and Billing

#### Billing & Quotas

* Tracks usage and applies pricing rules
* Enforces quotas, calculates cost
* Integrates with user tokens or credits

#### Audit & Logging

* Records all significant actions
* Optional tamper-resistant logging (hash-chaining)
* Supports forensic analysis and compliance

#### Monitoring / Observability

* Prometheus for metrics
* Jaeger (OpenTelemetry) for tracing
* Loki/FluentBit for centralized logs

---

## Storage Technologies

| Type      | Tool       | Purpose                        |
|-----------|------------|--------------------------------|
| SQL       | PostgreSQL | Tasks, users, billing metadata |
| In-memory | Redis      | Cache, heartbeats, state       |
| Object    | MinIO (S3) | Inputs, outputs, artifacts     |

---

## Communication Patterns

* **gRPC**: Internal service-to-service RPC
* **Kafka**: Async event streaming between components
* **REST**: External API for client/task submission

---

## Workflows

### Task Lifecycle (Simplified Diagram)

```text
User → API Gateway → Task Service
                       ↓
              Kafka: TaskCreated
                       ↓
                 Scheduler
                       ↓
          gRPC → Worker Manager
                       ↓
              Kafka: TaskAssigned
                       ↓
                   Worker
                       ↓
        ┌──────────────┬───────────────┐
        ↓                             ↓
 MinIO (partial output)     gRPC → Result Aggregator
                                      ↓
                        Kafka: SubTaskCompleted
                                      ↓
                           ResultAggregator
                                      ↓
                          Kafka: TaskCompleted
                                      ↓
                ┌──────────────┬───────────────┐
                ↓                              ↓
        Task Service                   Billing Service
```

---

## Protocol Buffers & Code Generation

The system uses Protocol Buffers to define service interfaces and messages. All proto files are located in the `proto/`
directory:

```
proto/
├── task/task.proto
├── scheduler/scheduler.proto
├── worker/worker.proto
├── result/result.proto
├── billing/billing.proto
├── audit/audit.proto
├── user/user.proto
└── kafka/messages.proto
```

To generate Go code:

```bash
./scripts/generate.sh
```

Outputs will be placed in `pkg/proto`, and include:

* gRPC clients and servers
* Typed message structs
* Efficient serialization logic

---

## Summary

This platform is purpose-built for scalable, distributed Go task execution — combining modern infrastructure tooling.

---
