package protoc

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

// Version returns the version of the protoc compiler.
// Taken from: https://github.com/grpc/grpc-go/blob/master/cmd/protoc-gen-go-grpc/grpc.go#L130
func Version(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}
