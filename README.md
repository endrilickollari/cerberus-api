# Cerberus API

A modern Go API for remote server management via SSH connections.


[![Swagger API Docs](https://img.shields.io/badge/API_Docs-Swagger-85ea2d?logo=swagger)](https://cerberus-api-0773eaec6d0f.herokuapp.com/swagger/index.html)

Explore the API endpoints and test requests directly via Swagger UI.

## Project Overview

Cerberus API allows you to securely connect to remote servers via SSH and retrieve system information such as server details, CPU information, disk usage, running processes, and Docker containers.

## Architecture

This project follows Clean Architecture principles to ensure scalability, maintainability, and testability:

```
┌───────────────────────────────────────────────────────────┐
│                     API Endpoints                         │
└─────────────────────────┬─────────────────────────────────┘
                          │
┌─────────────────────────▼─────────────────────────────────┐
│                      Handlers                             │
└─────────────────────────┬─────────────────────────────────┘
                          │
┌─────────────────────────▼─────────────────────────────────┐
│                      Services                             │
└─────────────────────────┬─────────────────────────────────┘
                          │
┌─────────────────────────▼─────────────────────────────────┐
│                   Repositories                            │
└─────────────────────────┬─────────────────────────────────┘
                          │
┌─────────────────────────▼─────────────────────────────────┐
│                  Infrastructure                           │
└───────────────────────────────────────────────────────────┘
```

### Layers

- **API Endpoints**: HTTP routes that accept requests and return responses
- **Handlers**: Process HTTP requests, call services, and format responses
- **Services**: Implement business logic and orchestrate repositories
- **Repositories**: Data access layer that abstracts storage mechanisms
- **Infrastructure**: Technical implementation details (SSH client, token management, etc.)

## Key Features

- **Clean Architecture**: Separation of concerns with distinct layers
- **Dependency Injection**: All dependencies are injected for better testability
- **Middleware Support**: Authentication middleware for protected routes
- **JWT Authentication**: Secure authentication with JWT tokens
- **Graceful Shutdown**: Proper server shutdown with timeout
- **Environment Configuration**: Configuration via environment variables
- **Swagger Documentation**: API documentation with Swagger
- **Error Handling**: Consistent error handling and responses

## API Endpoints

### Authentication

- `POST /login`: Authenticate with SSH credentials and receive a JWT token

### Server Details

- `GET /server-details`: Get basic server information
- `GET /server-details/cpu-info`: Get CPU information
- `GET /server-details/disk-usage`: Get disk usage information
- `GET /server-details/running-processes`: Get running processes information

### Docker

- `GET /docker/container-details`: Get information about Docker containers

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Access to a remote server with SSH enabled

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/cerberus-api.git
cd cerberus-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure environment variables:
```bash
export PORT=8080
export JWT_SECRET=your_secret_key
```

4. Run the application:
```bash
go run cmd/server/main.go
```

### Usage

1. Authenticate with SSH credentials:
```bash
curl -X POST http://localhost:8080/login -d '{
  "ip": "your-server-ip",
  "username": "your-username",
  "port": "22",
  "password": "your-password"
}'
```

2. Use the returned token for subsequent requests:
```bash
curl -H "Authorization: Bearer your-token" http://localhost:8080/server-details
```

## Code Structure

```
remote-server-api/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── config/
│   └── config.go                   # Configuration handling
├── internal/
│   ├── api/                        # API layer
│   │   ├── handlers/               # Request handlers
│   │   ├── router/                 # Router setup
│   │   ├── server/                 # HTTP server setup
│   │   └── response/               # Response handling
│   ├── domain/                     # Business domain
│   │   ├── auth/                   # Authentication domain
│   │   ├── server/                 # Server details domain
│   │   └── docker/                 # Docker domain
│   ├── infrastructure/             # Infrastructure concerns
│   │   ├── ssh/                    # SSH client
│   │   ├── persistence/            # Data persistence
│   │   └── token/                  # Token management
│   └── utils/                      # Utility functions
└── pkg/                            # Public packages
```

## Improvements from Previous Version

1. **Architecture**: Moved from a monolithic structure to a clean, layered architecture
2. **Router**: Replaced standard library's HTTP handler with Chi router for middleware support
3. **Authentication**: Centralized JWT validation in a middleware
4. **Error Handling**: Standardized error responses across the API
5. **Dependency Injection**: Implemented dependency injection for better testability
6. **Configuration**: Moved from hardcoded values to environment variables
7. **HTTP Methods**: Changed from POST to GET for data retrieval endpoints

## Security Considerations

- The JWT secret key should be stored securely and rotated regularly
- In production, consider implementing more robust error handling and logging
- For improved security, consider using SSH keys instead of passwords
- The current implementation uses `ssh.InsecureIgnoreHostKey()` which is not recommended for production

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request