package framework

import (
	"strings"

	"github.com/petergtz/pegomock/mockgen"
	"github.com/petergtz/pegomock/model"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/lovoo/protoc-gen-go-grpcmock/internal/generator"
)

const metadataPackage = protogen.GoImportPath("google.golang.org/grpc/metadata")

type pegomockMocker struct{}

func NewPegomockMocker() generator.Mocker {
	return &pegomockMocker{}
}

func (pm *pegomockMocker) Name() string {
	return "pegomock"
}

func (pm *pegomockMocker) Mock(g *protogen.GeneratedFile, file *protogen.File) {
	matchers := make(map[string]string)

	for _, service := range file.Services {
		pkg := string(file.GoPackageName)

		interfaces := []*model.Interface{
			{
				Name:    service.GoName + ClientSuffix, // Pegomock automatically adds the `Mock` prefix
				Methods: mapSlice(service.Methods, pm.clientMethod),
			},
			{
				Name:    service.GoName + ServerSuffix, // Pegomock automatically adds the `Mock` prefix
				Methods: mapSlice(service.Methods, pm.serverMethod),
			},
		}

		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				interfaces = append(interfaces, pm.generateClientStreamHandler(method))
				interfaces = append(interfaces, pm.generateServerStreamHandler(method))
			}
		}

		ast := &model.Package{
			Name:       pkg,
			Interfaces: interfaces,
		}

		data, types := mockgen.GenerateOutput(ast, file.Desc.Path(), "", pkg, string(file.GoImportPath))

		for t, matcher := range types {
			matchers[t] = strings.ReplaceAll(matcher, pkg+".", "")
		}

		// Strip the header comment and package name.
		// Package names are kept.
		g.P(substringAfter(string(data), "package "+pkg))
	}

	for t, matcher := range matchers {
		// The types context.Context and grpc.* must be excluded, since
		// they are not unique to a .proto file.
		if t == "context_context" || strings.HasPrefix(t, "grpc_") {
			continue
		}

		// Strip the header comment, package name and imports.
		g.P(substringAfter(matcher, ")"))
	}
}

func (pm *pegomockMocker) clientMethod(method *protogen.Method) *model.Method {
	m := &model.Method{
		Name: method.GoName,
		In: []*model.Parameter{
			{
				Name: "ctx",
				Type: &model.NamedType{
					Package: string(contextPackage),
					Type:    "Context",
				},
			},
		},
		Out: make([]*model.Parameter, 2), //nolint:gomnd
		Variadic: &model.Parameter{
			Name: "opts",
			Type: &model.NamedType{
				Package: string(grpcPackage),
				Type:    "CallOption",
			},
		},
	}

	if !method.Desc.IsStreamingClient() {
		m.In = append(m.In, &model.Parameter{
			Name: "in",
			Type: &model.PointerType{
				Type: pm.qualifiedGoIdent(method.Input.GoIdent),
			},
		})
	}

	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		m.Out[0] = &model.Parameter{
			Type: &model.PointerType{
				Type: pm.qualifiedGoIdent(method.Output.GoIdent),
			},
		}
	} else {
		m.Out[0] = &model.Parameter{
			Type: &model.NamedType{
				Package: string(method.Output.GoIdent.GoImportPath),
				Type:    method.Parent.GoName + "_" + method.GoName + ClientSuffix,
			},
		}
	}

	m.Out[1] = &model.Parameter{
		Type: model.PredeclaredType("error"),
	}

	return m
}

func (pm *pegomockMocker) serverMethod(method *protogen.Method) *model.Method {
	m := &model.Method{
		Name: method.GoName,
	}

	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		m.In = append(m.In, &model.Parameter{
			Name: "ctx",
			Type: &model.NamedType{
				Package: string(contextPackage),
				Type:    "Context",
			},
		})

		m.Out = append(m.Out, &model.Parameter{
			Type: &model.PointerType{
				Type: pm.qualifiedGoIdent(method.Output.GoIdent),
			},
		})
	}
	if !method.Desc.IsStreamingClient() {
		m.In = append(m.In, &model.Parameter{
			Name: "in",
			Type: &model.PointerType{
				Type: pm.qualifiedGoIdent(method.Input.GoIdent),
			},
		})
	}
	if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
		m.In = append(m.In, &model.Parameter{
			Name: "out",
			Type: &model.NamedType{
				Package: string(method.Output.GoIdent.GoImportPath),
				Type:    method.Parent.GoName + "_" + method.GoName + ServerSuffix,
			},
		})
	}

	m.Out = append(m.Out, &model.Parameter{
		Type: model.PredeclaredType("error"),
	})

	return m
}

func (pm *pegomockMocker) generateClientStreamHandler(method *protogen.Method) *model.Interface {
	i := &model.Interface{
		Name: method.Parent.GoName + "_" + method.GoName + ClientSuffix, // Pegomock automatically adds the `Mock` prefix
		Methods: []*model.Method{
			{
				Name: "Header",
				Out: []*model.Parameter{
					{
						Type: pm.qualifiedGoIdent(metadataPackage.Ident("MD")),
					},
					{
						Type: model.PredeclaredType("error"),
					},
				},
			},
			{
				Name: "Trailer",
				Out: []*model.Parameter{
					{
						Type: pm.qualifiedGoIdent(metadataPackage.Ident("MD")),
					},
				},
			},
			{
				Name: "CloseSend",
				Out: []*model.Parameter{
					{
						Type: model.PredeclaredType("error"),
					},
				},
			},
			{
				Name: "Context",
				Out: []*model.Parameter{
					{
						Type: pm.qualifiedGoIdent(contextPackage.Ident("Context")),
					},
				},
			},
			{
				Name: "SendMsg",
				In: []*model.Parameter{
					{
						Name: "msg",
						Type: model.PredeclaredType("interface{}"),
					},
				},
				Out: []*model.Parameter{
					{
						Type: model.PredeclaredType("error"),
					},
				},
			},
			{
				Name: "RecvMsg",
				In: []*model.Parameter{
					{
						Name: "msg",
						Type: model.PredeclaredType("interface{}"),
					},
				},
				Out: []*model.Parameter{
					{
						Type: model.PredeclaredType("error"),
					},
				},
			},
		},
	}

	if method.Desc.IsStreamingClient() {
		i.Methods = append(i.Methods, &model.Method{
			Name: "Send",
			In: []*model.Parameter{
				{
					Name: "m",
					Type: &model.PointerType{
						Type: pm.qualifiedGoIdent(method.Input.GoIdent),
					},
				},
			},
			Out: []*model.Parameter{
				{
					Type: model.PredeclaredType("error"),
				},
			},
		})
	}

	methodName := "Recv"
	if !method.Desc.IsStreamingServer() {
		methodName = "CloseAndRecv"
	}

	i.Methods = append(i.Methods, &model.Method{
		Name: methodName,
		Out: []*model.Parameter{
			{
				Type: &model.PointerType{
					Type: pm.qualifiedGoIdent(method.Output.GoIdent),
				},
			},
			{
				Type: model.PredeclaredType("error"),
			},
		},
	})

	return i
}

func (pm *pegomockMocker) generateServerStreamHandler(method *protogen.Method) *model.Interface {
	i := &model.Interface{
		Name: method.Parent.GoName + "_" + method.GoName + ServerSuffix, // Pegomock automatically adds the `Mock` prefix
		Methods: []*model.Method{
			{
				Name: "SetHeader",
				In: []*model.Parameter{
					{
						Name: "md",
						Type: pm.qualifiedGoIdent(metadataPackage.Ident("MD")),
					},
				},
				Out: []*model.Parameter{
					{
						Type: model.PredeclaredType("error"),
					},
				},
			},
			{
				Name: "SendHeader",
				In: []*model.Parameter{
					{
						Name: "md",
						Type: pm.qualifiedGoIdent(metadataPackage.Ident("MD")),
					},
				},
				Out: []*model.Parameter{
					{
						Type: model.PredeclaredType("error"),
					},
				},
			},
			{
				Name: "SetTrailer",
				In: []*model.Parameter{
					{
						Name: "md",
						Type: pm.qualifiedGoIdent(metadataPackage.Ident("MD")),
					},
				},
			},
			{
				Name: "Context",
				Out: []*model.Parameter{
					{
						Type: pm.qualifiedGoIdent(contextPackage.Ident("Context")),
					},
				},
			},
			{
				Name: "SendMsg",
				In: []*model.Parameter{
					{
						Name: "m",
						Type: model.PredeclaredType("interface{}"),
					},
				},
				Out: []*model.Parameter{
					{
						Type: model.PredeclaredType("error"),
					},
				},
			},
			{
				Name: "RecvMsg",
				In: []*model.Parameter{
					{
						Name: "m",
						Type: model.PredeclaredType("interface{}"),
					},
				},
				Out: []*model.Parameter{
					{
						Type: model.PredeclaredType("error"),
					},
				},
			},
		},
	}

	if method.Desc.IsStreamingClient() {
		i.Methods = append(i.Methods, &model.Method{
			Name: "Recv",
			Out: []*model.Parameter{
				{
					Type: &model.PointerType{
						Type: pm.qualifiedGoIdent(method.Input.GoIdent),
					},
				},
				{
					Type: model.PredeclaredType("error"),
				},
			},
		})
	}

	methodName := "Send"
	if !method.Desc.IsStreamingServer() {
		methodName = "SendAndClose"
	}

	i.Methods = append(i.Methods, &model.Method{
		Name: methodName,
		In: []*model.Parameter{
			{
				Name: "m",
				Type: &model.PointerType{
					Type: pm.qualifiedGoIdent(method.Output.GoIdent),
				},
			},
		},
		Out: []*model.Parameter{
			{
				Type: model.PredeclaredType("error"),
			},
		},
	})

	return i
}

func (pm *pegomockMocker) qualifiedGoIdent(ident protogen.GoIdent) model.Type {
	return &model.NamedType{
		Package: string(ident.GoImportPath),
		Type:    ident.GoName,
	}
}

func init() {
	setMocker("pegomock", NewPegomockMocker)
}
