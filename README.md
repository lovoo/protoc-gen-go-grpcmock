<p align="center">
    <h1>protoc-gen-go-grpcmock</h1>
    <span>Google protocol buffer compiler plugin to generate Mocks for gRPC Services in Go.</span>
</p>

##  Installation

Download the latest version from the [Release Page](https://github.com/lovoo/protoc-gen-go-grpcmock/releases/latest).
Extract the archive and make sure, the `protoc-gen-go-grpcmock` binary can be found in your `PATH`. 

For instance:
```
$ VERSION=$(curl -fsSL https://github.com/lovoo/protoc-gen-go-grpcmock/releases/latest -H "Accept: application/json"  | jq -r .tag_name)
$ curl -fsSL "https://github.com/lovoo/protoc-gen-go-grpcmock/releases/download/${VERSION}/protoc-gen-go-grpcmock_${VERSION:1}_$(uname -s)_$(uname -m).tar.gz" | tar -xzC /usr/local/bin protoc-gen-go-grpcmock
```

Or build the `protoc-gen-go-grpcmock` binary from source (requires Go 1.21+).

```
$ git clone https://github.com/lovoo/protoc-gen-go-grpcmock && cd protoc-gen-go-grpcmock
$ go build -ldflags "-X main.Version=$(git describe --tags)" cmd/protoc-gen-go-grpcmock
```

## Usage

Generate code by specifying the `--go-grpcmock_out` (and optional `--go-grpcmock_opt`) argument when invoking the `protoc` compiler.

```sh
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --go-grpcmock_out=. --go-grpcmock_opt=paths=source_relative \
    examples/helloworld/helloworld.proto
```

This will generate a `*_grpc_mock.pb` file for each specified `.proto` file, containing:

* Generated Client and Server Mocks for each Service
* Matchers for all Messages

## Options

The following parameters can be provided to change the behaviour of the compiler plugin.

| Parameter        | Default   | Available Options     | Description                   |
|------------------|-----------|-----------------------|-------------------------------|
| `framework`      | "testify" | "testify", "pegomock" | The mocking framework to use. |
| `import_package` | false     | true/false            | Import the file's Go package. <br /> This can be useful if mocks should be generated <br /> in a different package, then the original `.pb.go` files |

## Examples

Examples can be found in the [examples](./examples) directory.

## License

The MIT License (MIT). Please see [LICENSE](LICENSE) for more information.
