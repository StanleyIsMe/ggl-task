# Gogolook Task API

In this project, I implement a simple API for Gogolook Interview Task.

## Tech Stack

- Go 1.23+
- Docker
- Zerolog for structured logging
- Comprehensive test coverage with race/leak detection
- Trivy/Govulncheck for security scanning
- Makefile for build, test, lint, etc.
- Github Actions for CI

## Prerequisites

- Go 1.23+
- Docker
- Make

## Getting Started

1. Clone the repository:
    ```sh
    git clone https://github.com/StanleyIsMe/ggl-task.git
    ```

2. Clone the `config/api/base.yaml` file to `config/api/local.yaml` and fill in the details:
    ```sh
    cp config/api/base.yaml config/api/local.yaml
    ```

3. Build the Docker image:
    ```sh
    make build
    ```

4. Run the Docker image:
    ```sh
    make up
    ```

5. Run your app, and browse to http://localhost:8080/swagger/index.html. You will see Swagger 2.0 Api documents.

## Directory Structure

```sh
.
├── cmd                    # application entry point, usually a main.go file
│   └── api
├── config                 # configuration files
│   └── api
├── database               # database related files, including migrations
│   ├── init
│   └── migrations
├── docs                   # swagger api documentation
├── internal               # core logic
│   ├── api                # initialize api server, for example, http server, grpc server, etc.
│   │   ├── config         # configuration for api server
│   │   ├── middleware
│   │   ├── repository
│   │   └── server
│   └── task                 # domain logic. naming is depends on the business
│       ├── delivery         # delivery layer is responsible for handling http/grpc request and response
│       │   └── http
│       ├── domain           # domain layer is responsible for defining the business logic
│       │   ├── entities
│       │   ├── mock
│       │   ├── repository
│       │   └── usecase
│       ├── mock               # mock files for unit test
│       │   ├── repositorymock
│       │   └── usecasemock
│       ├── repository         # implementing the data/external service access logic 
│       │   └── memory
│       └── usecase            # implementing the business logic 
└── pkg                        # internal packages
    ├── config
    ├── shutdown
    └── transport
        └── middleware
```

