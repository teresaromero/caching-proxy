# Caching Proxy

[![compile and test](https://github.com/teresaromero/caching-proxy/actions/workflows/compile-and-test.yml/badge.svg)](https://github.com/teresaromero/caching-proxy/actions/workflows/compile-and-test.yml)

[![golangci-lint](https://github.com/teresaromero/caching-proxy/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/teresaromero/caching-proxy/actions/workflows/golangci-lint.yml)

A command-line interface (CLI) tool for a caching server, built using Go. This project aims to provide a simple and efficient caching mechanism to improve the performance of web applications by reducing the load on the backend servers.

Originally posted at [roadmap.sh](https://roadmap.sh/projects/caching-server)

## Features

- **In-memory caching**: Store frequently accessed data in memory to reduce latency.
- **TTL (Time-to-Live)**: Automatically expire cache entries after a specified duration.
- **LRU (Least Recently Used)**: Evict the least recently used items when the cache reaches its capacity.
- **CLI Interface**: Easy-to-use command-line interface for managing the cache.

## Installation

To install the Caching Proxy CLI, ensure you have Go installed, then run:

```sh
go install github.com/teresaromero/caching-proxy@latest
```

## Usage

### Starting the Server

Start the caching server with the following command:

```sh
caching-proxy --port <number> --origin <url>
```

### Clear cache

Delete data from the cache using:

```sh
caching-proxy --clear-cache
```

## Configuration

You can configure the caching server using a configuration file or environment variables. The default configuration file is `config.yaml`.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or suggestions, please open an issue.
