# DNS Server in Go

This project implements a simple DNS server in Go, located in `app/main.go`. It listens for DNS queries over UDP and responds with a hardcoded answer for A record queries (IPv4 address for 127.0.0.1).

## Features

- Listens on UDP port 2053 (127.0.0.1:2053)
- Parses DNS queries and constructs valid DNS responses
- Responds to A record (IPv4) queries with 127.0.0.1
- Handles basic DNS header fields and question/answer sections
- Ignores unsupported query types and classes

## Getting Started

### Prerequisites

- Go 1.24 or higher

### Installation

Clone the repository:

```bash
git clone https://github.com/ronitrajfr/dns-server.git
cd dns-server
```

### Running the Server

```bash
go run app/main.go
```

The server will start and listen for UDP DNS queries on 127.0.0.1:2053.

### Testing the Server

You can test the server using `dig`:

```bash
dig @127.0.0.1 -p 2053 example.com A
```

You should receive a response with the IP address 127.0.0.1 for any A record query.

## Code Overview

- `app/main.go`: Main server implementation. Handles UDP socket, parses DNS packets, and constructs responses.
- `go.mod`: Go module definition.

## How It Works

- The server listens for UDP packets on 127.0.0.1:2053.
- For each incoming DNS query:
  - Parses the DNS header and question section.
  - If the query is for an A record (QTYPE=1, QCLASS=1), responds with 127.0.0.1.
  - For unsupported query types or classes, no response is sent.

## Development

- Modify `app/main.go` to extend functionality (e.g., support more record types, logging, etc.).
- Use print statements for debugging; output will appear in the console.

## License

This project is for educational purposes and is based on the Codecrafters DNS server challenge.
