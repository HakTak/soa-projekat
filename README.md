# SOA Tourism Application

A comprehensive microservices-based tourism platform built with a service-oriented architecture. The application enables tourists to discover, plan, and execute guided tours, while guides can create and manage their offerings.

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Microservices](#microservices)
- [Features](#features)
- [User Roles](#user-roles)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
- [API Gateway](#api-gateway)
- [Development Checkpoints](#development-checkpoints)
- [Contributing](#contributing)

## Project Overview

This is a distributed tourism application built following microservices architecture principles. The system handles user authentication, tour management, blog functionality, follower relationships, and tour execution with real-time position tracking.

Key architectural principles:
- **Microservices**: Independent, loosely coupled services
- **Database per Service**: Each microservice manages its own data
- **API Gateway**: Single entry point for all client requests
- **Containerization**: Docker-based deployment
- **NoSQL & Relational Databases**: Mixed database strategy (MongoDB, Neo4j, PostgreSQL)
- **gRPC Communication**: Efficient inter-service communication
- **Distributed Tracing & Logging**: Centralized observability

## Architecture

```
┌─────────────────┐
│   Frontend      │
│    (Angular)    │
└────────┬────────┘
         │
    ┌────▼────┐
    │ Gateway  │
    └────┬────┘
         │
    ┌────┴────────────────┬──────────────┬──────────────┐
    │                     │              │              │
┌───▼───┐        ┌────────▼──┐   ┌──────▼──┐   ┌──────▼───┐
│  Auth │        │ Stakeholder │  │  Blog   │   │  Tour    │
│       │        │   Service   │  │ Service │   │ Service  │
└───────┘        └─────────────┘  └─────────┘   └──────────┘
                                          │
                                    ┌─────▼──────┐
                                    │ Follower   │
                                    │ Service    │
                                    └────────────┘
```

## Technology Stack

### Backend Services
- **Go**: API Gateway, Auth, Follower, Shopping Cart, Stakeholders, Tour services
- **C# / .NET**: Blog service
- **Protocol Buffers**: gRPC service definitions
- **Python**: Potential utility services

### Databases
- **MongoDB**: NoSQL document storage (Blog service, main data)
- **Neo4j**: Graph database (Follower relationships and recommendations)
- **PostgreSQL/SQL**: Relational data storage

### Infrastructure & DevOps
- **Docker**: Container orchestration
- **Docker Compose**: Multi-container application management
- **gRPC**: Efficient inter-service communication
- **REST API**: Client-service communication

### Frontend
- **Angular**: Modern web framework
- **TypeScript**: Type-safe JavaScript
- **Maps Integration**: Tour mapping and position tracking

### Observability
- **Tracing**: Distributed request tracing
- **Logging**: Centralized log aggregation (Promtail/Loki)
- **Metrics**: System and container monitoring
- **Visualization**: Grafana dashboards

## Microservices

### 1. API Gateway
**Location**: `API_GATEWAY/`  
**Language**: Go  
**Purpose**: Single entry point for all client requests
- Request routing to appropriate microservices
- Authentication middleware
- CORS handling
- Distributed tracing

**Port**: 8000 (default)

### 2. Auth Service
**Location**: `AUTH/`  
**Language**: Go  
**Purpose**: User authentication and authorization
- User registration
- Login/logout functionality
- Token generation and validation
- Role-based access control

**Port**: 8001 (default)

### 3. Stakeholders Service
**Location**: `STAKEHOLDERS/`  
**Language**: Go  
**Purpose**: User profile and admin operations
- User profile management
- User account information
- Admin operations (view users, block accounts)
- Profile updates (name, bio, picture)

**Port**: 8002 (default)

### 4. Blog Service
**Location**: `BLOG/`  
**Language**: C# / .NET  
**Purpose**: Blog and comment management
- Blog post creation and management
- Comments on blog posts
- Blog post liking/unliking
- Markdown support for blog descriptions

**Port**: 8003 (default)  
**Database**: MongoDB

### 5. Follower Service
**Location**: `FOLLOWER/`  
**Language**: Go  
**Purpose**: User relationship and recommendation engine
- Follow/unfollow users
- Access control based on follow relationships
- Tour recommendations
- Graph-based relationship management

**Port**: 8004 (default)  
**Database**: Neo4j (Graph Database)

### 6. Tour Service
**Location**: `TOUR/`  
**Language**: Go  
**Purpose**: Tour management and execution
- Tour creation and publishing
- Tour execution tracking
- Waypoint management
- Position simulation

**Port**: 8005 (default)

### 7. Shopping Cart Service
**Location**: `SHOPPING-CART/`  
**Language**: Go  
**Purpose**: Shopping cart and purchase management
- Add/remove tours from cart
- Cart total calculation
- Checkout and token generation
- Purchase history

**Port**: 8006 (default)

### 8. Common/Shared
**Location**: `COMMON/`  
**Purpose**: Shared resources across services
- Protocol Buffer definitions
- Shared utilities and middleware
- Authentication utilities
- Metadata extraction

## Features

### User Authentication & Authorization
- ✅ User registration with role selection (Guide, Tourist, Admin)
- ✅ Admin user management and account blocking
- ✅ Profile management and updates
- ✅ Role-based access control

### Blog System
- ✅ Create blog posts with markdown support
- ✅ Add comments to blog posts
- ✅ Like/unlike functionality with count tracking
- ✅ Restricted access to followed users' blogs

### Tour Management
- ✅ Create and manage tours with draft/published/archived states
- ✅ Define waypoints with geographic coordinates
- ✅ Multiple transport type support (walking, cycling, driving)
- ✅ Automatic distance calculation between waypoints
- ✅ Tour reviews and ratings
- ✅ Tour discovery and browsing

### Social Features
- ✅ Follow/unfollow other users
- ✅ Follower-based access control
- ✅ User recommendations based on graph relationships
- ✅ Feed based on followed users' content

### Shopping & Purchases
- ✅ Add tours to shopping cart
- ✅ Checkout and purchase confirmation
- ✅ Purchase tokens for tour access
- ✅ Order history

### Tour Execution
- ✅ Start/complete/abandon tours
- ✅ Real-time position tracking (position simulator)
- ✅ Waypoint completion tracking
- ✅ Activity timestamp recording
- ✅ Session management

### Advanced Features
- ✅ gRPC communication between services
- ✅ SAGA pattern for distributed transactions
- ✅ Distributed tracing and logging
- ✅ System and container metrics
- ✅ Centralized gateway routing

## User Roles

### Administrator
- View all user accounts
- Block user accounts
- System administration

### Guide (Vodič)
- Create and manage tours
- Manage tour waypoints
- Publish/archive tours
- View tour reviews and ratings

### Tourist (Turista)
- Browse published tours
- Follow guides and other tourists
- Purchase tours
- Execute tours with position tracking
- Leave reviews and ratings
- Create blog posts
- View followed users' content

## Prerequisites

### Required Software
- **Docker**: >= 20.10
- **Docker Compose**: >= 1.29
- **Go**: >= 1.16 (for Go services)
- **Node.js**: >= 14 (for frontend)
- **.NET SDK**: >= 6.0 (for Blog service)
- **Git**: >= 2.30

### Optional but Recommended
- **MongoDB**: For development/testing
- **Neo4j**: For testing graph features
- **Postman**: For API testing
- **Git Bash**: For Windows users

## Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd "4. god/SOA/PROJEKAT"
```

### 2. Configure Environment Variables

Create a `.env` file in the project root:
```env
# Database Configuration
MONGODB_URI=mongodb://mongodb:27017
NEO4J_URI=neo4j://neo4j:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=password
POSTGRES_URL=postgres://user:password@postgres:5432/tourism_db

# Service Configuration
JWT_SECRET=your-secret-key-here
API_GATEWAY_PORT=8000
AUTH_SERVICE_PORT=8001
STAKEHOLDERS_PORT=8002
BLOG_SERVICE_PORT=8003
FOLLOWER_SERVICE_PORT=8004
TOUR_SERVICE_PORT=8005
SHOPPING_CART_PORT=8006

# Frontend
FRONTEND_PORT=4200
FRONTEND_API_URL=http://localhost:8000
```

### 3. Build Services Locally (Optional)

#### Go Services
```bash
# For each Go service (replace SERVICE_NAME with directory name)
cd SERVICE_NAME
go mod download
go build -o ./bin/app ./cmd/main.go
cd ..
```

#### .NET Blog Service
```bash
cd BLOG
dotnet restore
dotnet build
cd ..
```

#### Frontend
```bash
cd FRONT
npm install
cd ..
```

## Running the Application

### Using Docker Compose (Recommended)

1. **Build and start all services**:
   ```bash
   docker-compose up -d
   ```

2. **Verify services are running**:
   ```bash
   docker-compose ps
   ```

3. **View logs**:
   ```bash
   docker-compose logs -f [service-name]
   ```

4. **Stop services**:
   ```bash
   docker-compose down
   ```

### Using Individual Docker Containers

Each service has a `Dockerfile` in its directory. Build and run individually:

```bash
# Build
docker build -t tourism-service-name ./SERVICE_NAME

# Run
docker run -p PORT:PORT --network tourism-network tourism-service-name
```

### Local Development (Without Docker)

1. **Start Auth Service**:
   ```bash
   cd AUTH
   go run ./cmd/main.go
   ```

2. **Start API Gateway** (in new terminal):
   ```bash
   cd API_GATEWAY
   go run ./cmd/main.go
   ```

3. **Start other services** similarly

4. **Start Frontend** (in new terminal):
   ```bash
   cd FRONT
   npm start
   ```

## Project Structure

```
PROJEKAT/
├── API_GATEWAY/           # Request routing and middleware
├── AUTH/                  # Authentication & authorization
├── BLOG/                  # Blog posts and comments (C#/.NET)
├── COMMON/                # Shared proto definitions and utilities
├── FOLLOWER/              # User relationships and recommendations
├── FRONT/                 # Angular frontend application
├── SHOPPING-CART/         # Shopping cart and purchases
├── STAKEHOLDERS/          # User profiles and admin operations
├── TOUR/                  # Tour management and execution
├── docker-compose.yml     # Multi-container orchestration
├── promtail-config.yml    # Logging configuration
└── README.md              # This file
```

## API Gateway

The API Gateway serves as the single entry point for all client requests. It handles:

- **Route Forwarding**: Maps client requests to appropriate microservices
- **Authentication**: Validates JWT tokens
- **CORS**: Cross-origin request handling
- **Tracing**: Distributed request tracking
- **Rate Limiting**: Request throttling (optional)

### Gateway Routes

```
GET    /auth/register              → Auth Service
POST   /auth/login                 → Auth Service
GET    /api/users/:id              → Stakeholders Service
PUT    /api/users/:id              → Stakeholders Service
GET    /api/blogs                  → Blog Service
POST   /api/blogs                  → Blog Service
POST   /api/tours                  → Tour Service
GET    /api/tours/:id              → Tour Service
POST   /api/tours/:id/execute      → Tour Service
GET    /api/cart                   → Shopping Cart Service
```

## Development Checkpoints

The project follows three evaluation checkpoints (KT):

### Checkpoint 1 (Week 21.07)
- Basic microservices (Auth, Stakeholders, Blog)
- User authentication and profiles
- Blog CRUD operations
- Docker containerization

**Required Points**: 3 per student

### Checkpoint 2 (Week 8.09)
- Follower service with Neo4j
- API Gateway implementation
- NoSQL integration (MongoDB)
- Docker Compose orchestration
- gRPC communication setup

**Required Points**: 3 per student

### Checkpoint 3 (Final Defense)
- RPC protocol implementation
- SAGA pattern for transactions
- Distributed tracing and logging
- System and container metrics
- Complete feature implementation

**Required Points**: 3 per student

## Communication Patterns

### REST API
- Client ↔ API Gateway
- Internal service discovery endpoints (optional)

### gRPC
- Gateway ↔ Microservices
- Inter-service communication
- Efficient binary protocol

### Event-Driven (SAGA)
- Distributed transactions across multiple services
- Choreography-based coordination
- Rollback handling

## Monitoring & Observability

### Distributed Tracing
- Request tracking across all services
- Visualization in tracing tools (Jaeger/Zipkin)
- Performance bottleneck identification

### Logging
- Centralized log aggregation
- Promtail → Loki stack
- Structured logging with request IDs

### Metrics
- Host machine metrics (CPU, RAM, network)
- Container metrics (resource usage)
- Application metrics (request rates, latencies)
- Grafana dashboards for visualization

## Deployment

### Production Deployment
```bash
# Using Docker Compose in production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Using Kubernetes (optional)
kubectl apply -f k8s/
```

### Health Checks
Each service exposes a health check endpoint:
```bash
curl http://localhost:8000/health
```

## API Testing

### Using Postman
- Import the `BLOG.http` file or equivalent API specifications
- Configure environment variables for different service endpoints
- Test individual endpoints and workflows

### Using cURL
```bash
# Register user
curl -X POST http://localhost:8000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"pass","email":"user@example.com","role":"tourist"}'

# Login
curl -X POST http://localhost:8000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"pass"}'
```

## Troubleshooting

### Services Won't Start
1. Check port availability: `netstat -an | grep LISTEN`
2. Verify Docker is running: `docker --version`
3. Check logs: `docker-compose logs service-name`

### Database Connection Issues
1. Verify MongoDB is running: `docker-compose ps mongodb`
2. Check connection strings in `.env` file
3. Ensure database containers are on the same network

### gRPC Communication Errors
1. Verify proto files are compiled
2. Check service addresses and ports
3. Ensure firewall rules allow inter-container communication

## Contributing

When contributing to this project:

1. Create a feature branch for each service/feature
2. Follow the existing code structure and conventions
3. Implement proper error handling and logging
4. Add appropriate Docker configuration
5. Test with Docker Compose before submitting
6. Update documentation as needed

## Contact & Support

For questions or issues related to this project, please contact the development team or create an issue in the repository.

---

**Last Updated**: May 2026  
**Version**: 1.0  
**Status**: Development Phase
