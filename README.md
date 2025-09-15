## gRPC Authentication Service
This is a standalone microservice written in Go that provides user registration and login functionality using gRPC. It's designed as a foundational service that other applications could use for authentication.

The service uses a SQLite database to store user information and follows security best practices by hashing passwords with bcrypt.

### Features
- User Registration: Securely register new users. Passwords are never stored in plain text.
- User Login: Authenticate users and issue a JWT access token upon successful login.
- gRPC Interface: All communication is handled through a clearly defined gRPC protocol.
- SQLite Persistence: User data is stored in a simple, file-based SQLite database.
- Docker: Comes with a `Dockerfile` and `docker-compose.yml` for easy setup.

### Prerequisites
- Go
- Docker
- protoc & Go Plugins (Required for local development if you modify the .proto file)

### Running Locally

1. Start the Service: `docker-compose up --build -d`
2. Use the client
```
# Register a new user
go run ./client register testuser securepassword123

# Log in as the user
go run ./client login testuser securepassword123

# View registered users
go run ./client list
```

### Managing the Service
```
# Viewing the logs
docker-compose logs -f

# Stopping the service
docker-compose down
```
