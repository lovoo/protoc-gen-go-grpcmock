package framework

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/lovoo/protoc-gen-go-grpcmock/internal/generator"
)

const (
	MockPrefix   = "Mock"
	ClientSuffix = "Client"
	ServerSuffix = "Server"

	contextPackage = protogen.GoImportPath("context")
	grpcPackage    = protogen.GoImportPath("google.golang.org/grpc")
)

var mocker = make(map[string]func() generator.Mocker)

var errUnknownMocker = errors.New("protoc-gen-go-grpcmock: unknown test framework")

func Mocker(name string) (generator.Mocker, error) {
	m, ok := mocker[name]
	if !ok {
		return nil, fmt.Errorf("%w %q. Please use one of the following: [%s]", errUnknownMocker, name, availableMocker())
	}
	return m(), nil
}

func setMocker(name string, ctor func() generator.Mocker) {
	mocker[name] = ctor
}

func availableMocker() string {
	m := make([]string, 0, len(mocker))
	for k := range mocker {
		m = append(m, fmt.Sprintf("%q", k))
	}
	return strings.Join(m, ", ")
}
