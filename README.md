# Aztec P2P Dial Tools

A pair of command-line tools for querying Aztec blockchain nodes via their P2P protocols. These tools allow you to check connectivity and retrieve status information from Aztec nodes using their libp2p peer IDs and addresses.

## Project Structure

```
aztec-p2pdial-mwe/
├── cmd/
│   ├── ping/           # Ping command source
│   │   └── main.go
│   └── status/         # Status command source
│       └── main.go
├── go.mod
├── go.sum
└── README.md
```

## Tools Overview

### `ping` (formerly `aztec_ping`)
Tests connectivity to an Aztec node using the ping protocol (`/aztec/req/ping/0.1.0`). Returns a simple "pong" response if the node is reachable and responsive.

### `status` (formerly `aztec_status`)
Queries an Aztec node's status using the status protocol (`/aztec/req/status/0.1.0`). Returns detailed status information from the node.

## Prerequisites

- Go 1.21 or later
- Access to an Aztec node's P2P multiaddress

## Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd aztec-p2pdial-mwe
```

2. Build the tools:
```bash
# Build both tools
go build -o ping ./cmd/ping
go build -o status ./cmd/status

# Or build them individually
go build -o ping ./cmd/ping
go build -o status ./cmd/status
```

## Usage

Both tools require a target peer multiaddress in the format:
```
/ip4/<IP_ADDRESS>/tcp/40400/p2p/<PEER_ID>
```

### Ping Tool

Test connectivity to an Aztec node:

```bash
# Using built binary
./ping -d /ip4/<IP_ADDRESS>/tcp/40400/p2p/<PEER_ID>

# Or run directly with go
go run ./cmd/ping -d /ip4/<IP_ADDRESS>/tcp/40400/p2p/<PEER_ID>
```

**Expected Output:**
- `✅ pong` - Node is reachable and responding
- Error message - Node is unreachable or not responding

### Status Tool

Query detailed status from an Aztec node:

```bash
# Using built binary
./status -d /ip4/<IP_ADDRESS>/tcp/40400/p2p/<PEER_ID>

# Or run directly with go
go run ./cmd/status -d /ip4/<IP_ADDRESS>/tcp/40400/p2p/<PEER_ID>
```

**Expected Output:**
- Status response with detailed node information
- Response displayed in both string and hexadecimal format

## Command Line Options

Both tools support the following options:

- `-d <multiaddr>` - **Required**. The target peer multiaddress to connect to.

## Development Commands

### Building

```bash
# Build both commands
go build -o ping ./cmd/ping
go build -o status ./cmd/status

# Build with custom output location
go build -o bin/ping ./cmd/ping
go build -o bin/status ./cmd/status
```

### Running Without Building

```bash
# Run ping command directly
go run ./cmd/ping -d <multiaddr>

# Run status command directly
go run ./cmd/status -d <multiaddr>
```

### Clean Build

```bash
# Remove built binaries
rm -f ping status

# Or if you used custom output location
rm -rf bin/
```

## Technical Details

### Dependencies

- **libp2p**: Core P2P networking library
- **go-multiaddr**: Multiaddress parsing and handling
- **go-log**: Logging utilities (suppressed for cleaner output)

### Protocols Used

- **Ping Protocol**: `/aztec/req/ping/0.1.0`
- **Status Protocol**: `/aztec/req/status/0.1.0`

### Connection Details

- Uses RSA 2048-bit key pairs for identity
- TCP transport only
- No relay support (direct connections only)
- Temporary connections (tools exit after single request)

## Troubleshooting

### Common Issues

1. **"Please provide a target peer"**
   - Ensure you're using the `-d` flag with a valid multiaddress

2. **"failed to open stream"**
   - Check that the target IP and port are correct
   - Verify the peer ID matches the target node
   - Ensure the node is running and accepting connections

3. **"Invalid multiaddr"**
   - Verify the multiaddress format: `/ip4/<IP>/tcp/<PORT>/p2p/<PEER_ID>`
   - Check that the peer ID is a valid base58-encoded libp2p peer ID

### Network Requirements

- Direct network connectivity to the target node
- No firewall blocking the specified TCP port
- Target node must support the Aztec P2P protocols

## Development

### Building from Source

```bash
go mod download
go build -o ping ./cmd/ping
go build -o status ./cmd/status
```

### Testing

To test the tools, you need access to a running Aztec node. Replace the example addresses with actual node addresses:

```bash
# Test ping
./ping -d /ip4/<NODE_IP>/tcp/<NODE_PORT>/p2p/<NODE_PEER_ID>

# Test status
./status -d /ip4/<NODE_IP>/tcp/<NODE_PORT>/p2p/<NODE_PEER_ID>
```

## License

This project is provided as-is for testing and development purposes. 