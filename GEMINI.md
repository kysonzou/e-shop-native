
# GEMINI.md

## Project Overview

This project, `e-shop-native`, is a Go-based e-commerce backend. It follows a microservices architecture, with services like `user-srv` and `product-srv`. The services communicate using gRPC, and a gRPC-Gateway is used to expose a RESTful API.

The project uses the following key technologies:

*   **Go:** The primary programming language.
*   **gRPC:** For high-performance, cross-platform, inter-service communication.
*   **gRPC-Gateway:** To expose a RESTful JSON API from the gRPC services.
*   **Docker and Docker Compose:** For containerization and orchestration of the services and their dependencies.
*   **etcd:** For service discovery.
*   **MySQL:** As the primary data store.
*   **Redis:** For caching.
*   **Viper:** For configuration management.
*   **Zap:** For structured, leveled logging.
*   **GORM:** As the ORM for interacting with the MySQL database.
*   **Wire:** For compile-time dependency injection.

## Building and Running

The project uses a `Makefile` to simplify common development tasks.

### Key Commands

*   **Run the user service:**
    ```bash
    make run_user_srv
    ```

*   **Build the binary:**
    ```bash
    make build
    ```

*   **Run tests:**
    ```bash
    make test
    ```

*   **Generate Go code from `.proto` files:**
    ```bash
    make api
    ```

*   **Generate dependency injection code:**
    ```bash
    make wire
    ```

*   **Clean up build artifacts:**
    ```bash
    make clean
    ```

### Docker

The project includes a `docker-compose.yaml` file to run the required services (etcd, MySQL, Redis) in Docker containers.

*   **Start all services:**
    ```bash
    docker-compose up -d
    ```

*   **Stop all services:**
    ```bash
    docker-compose down
    ```

## Development Conventions

*   **Project Structure:** The project follows the standard Go project layout, with `cmd`, `internal`, and `pkg` directories.
*   **Configuration:** Configuration is managed through a `config.yaml` file and loaded using Viper.
*   **API:** The API is defined using Protocol Buffers in the `api/protobuf` directory.
*   **Dependency Injection:** The project uses `wire` for compile-time dependency injection.
*   **Testing:** The project uses the standard Go testing framework.
