<p align="center">
    <h1>protoc-gen-go-grpcmock</h1>
    <span>Google protocol buffer compiler plugin to generate Mocks for gRPC Services in Go.</span>
</p>

##  Installation

Install the latest version of the plugin with the `go install` command:

```
$ go install github.com/lovoo/protoc-gen-go-grpcmock/cmd/protoc-gen-go-grpcmock@latest
```

Or download the latest version from the [Release Page](https://github.com/lovoo/protoc-gen-go-grpcmock/releases/latest).

Make sure, the `protoc-gen-go-grpcmock` binary can be found in your `PATH`.

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