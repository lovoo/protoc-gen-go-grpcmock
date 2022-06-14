// protoc-gen-go-grpcmock is a plugin for the Google protocol buffer compiler to
// generate Go code. Install it by building this program and making it
// accessible within your PATH with the name:
//	protoc-gen-go-grpcmock
//
// The 'go-grpcmock' suffix becomes part of the argument for the protocol compiler,
// such that it can be invoked as:
//	protoc --go-grpcmock_out=. path/to/file.proto
//
// This generates Go service definitions for the protocol buffer defined by
// file.proto.  With that input, the output will be written to:
//	path/to/file_grpc_mock.pb.go
package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/lovoo/protoc-gen-go-grpcmock/internal/framework"
	"github.com/lovoo/protoc-gen-go-grpcmock/internal/generator"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("protoc-gen-go-grpcmockmock %v\n", version)
		return
	}

	var flags flag.FlagSet
	testFramework := flags.String("framework", "testify", "The mocking framework to use.")
	importPackage := flags.Bool("import_package", false, "Import the file's Go package.")
	protogen.Options{ParamFunc: flags.Set}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		m, err := framework.Mocker(*testFramework)
		if err != nil {
			return err
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			if *importPackage {
				// Force import of the package.
				f.GoImportPath = protogen.GoImportPath("")
			}

			generator.GenerateFile(version, gen, f, m)
		}

		return nil
	})
}
