## gRPC Authentication Service
This is a standalone microservice written in Go that provides user registration and login functionality using gRPC. It's designed as a foundational service that other applications could use for authentication.

The service uses a SQLite database to store user information and follows security best practices by hashing passwords with bcrypt.

### Features
- User Registration: Securely register new users. Passwords are never stored in plain text.
- User Login: Authenticate users and issue a JWT access token upon successful login.
- gRPC Interface: All communication is handled through a clearly defined gRPC protocol.
- SQLite Persistence: User data is stored in a simple, file-based SQLite database.

### Prerequisites
- Go (version 1.18 or later)
- Protocol Buffer Compiler (protoc) - Can be installed with brew install protobuf on macOS.
- Go plugins for protoc: `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28` & `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2`

### Running Locally

### Configuration
The server can be configured using the following environment variables. If they are not set, sensible defaults will be used.
- `JWT_SECRET`: The secret key used to sign JSON Web Tokens. (`export JWT_SECRET="my_secret"`)
- `DB_PATH`: The file path for the SQLite database. (Default: ./users.db)
- `PORT`: The port for the gRPC server to listen on. (Default: :50051)

1. Clone the repository: git clone `https://github.com/adammcgrogan/go-auth`, `cd go-auth`
2. Generate Code:
```protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    auth/auth.proto```
3. Start the Server: (In Terminal 1) `go run ./server`
4. Use the Client: (In Terminal 2) `go run ./client register testuser password123` & `go run ./client login testuser password123` 

