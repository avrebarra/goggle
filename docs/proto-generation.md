# Proto Generation

This project uses a parameterized script for generating gRPC code from Protocol Buffer files.

## Main Service Proto Generation

For the main Goggle service, use the Makefile target:

```bash
make generate-proto
```

This generates Go code from `proto/goggle.proto` into `proto/generated/`.

## Client Module Proto Generation

Client modules should be generated manually as needed by developers:

### Auth Client

```bash
./scripts/genproto -f internal/module/clientauth/proto/auth.proto -o internal/module/clientauth/proto/generated
```

### Other Client Modules

For any external service client, use the script directly:

```bash
# Example: Payment service client
./scripts/genproto -f proto/external/payment/payment.proto -o internal/module/clientpayment/proto/generated

# Example: User service client  
./scripts/genproto -f proto/external/user/user.proto -o internal/module/clientuser/proto/generated
```

## Script Options

```bash
./scripts/genproto --help
```

**Options:**
- `-f, --file PROTO_FILE`: Proto file to generate (required)
- `-o, --output OUTPUT_DIR`: Output directory for generated files (required)
- `-h, --help`: Show help message

## Requirements

- Docker must be running (uses `namely/protoc-all:1.51_1` image)
- Ruby for running the generation script
