package framework

import (
	"fmt"
	"strings"

	_ "github.com/stretchr/testify/mock" // needed for version information in the import path
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/lovoo/protoc-gen-go-grpcmock/internal/generator"
	"github.com/lovoo/protoc-gen-go-grpcmock/internal/model"
)

const (
	deprecationComment = "// Deprecated: Do not use."

	grpcMetaPackage    = protogen.GoImportPath("google.golang.org/grpc/metadata")
	testifyMockPackage = protogen.GoImportPath("github.com/stretchr/testify/mock")
)

type testifyMocker struct{}

func NewTestifyMocker() generator.Mocker {
	return &testifyMocker{}
}

func (tm *testifyMocker) Name() string {
	return "testify"
}

func (tm *testifyMocker) Mock(g *protogen.GeneratedFile, file *protogen.File) {
	for _, msg := range file.Messages {
		tm.generateMatcher(g, file.GoPackageName, msg.GoIdent.GoName)
	}

	for _, service := range file.Services {
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				tm.generateMatcher(g, file.GoPackageName, method.Parent.GoName+"_"+method.GoName+ClientSuffix)
				tm.generateMatcher(g, file.GoPackageName, method.Parent.GoName+"_"+method.GoName+ServerSuffix)
			}
		}

		tm.generateService(g, service)
	}
}

func (tm *testifyMocker) generateMatcher(g *protogen.GeneratedFile, pkg protogen.GoPackageName, typeName string) {
	g.P("func Any", typeName, "() ", testifyMockPackage.Ident("AnythingOfTypeArgument"), " {")
	g.P("return ", testifyMockPackage.Ident("AnythingOfType"), "(\"*", pkg, ".", typeName, "\")")
	g.P("}")
	g.P()
}

func (tm *testifyMocker) generateService(g *protogen.GeneratedFile, service *protogen.Service) {
	clientName := MockPrefix + service.GoName + ClientSuffix

	// Client structure.
	tm.generateStruct(g, clientName)

	// NewClient factory.
	tm.generateNewFunc(g, service, clientName)

	// Client method implementations.
	for _, method := range service.Methods {
		tm.generateMethodDefinitions(g, tm.clientMethod(g, method))
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			tm.generateClientStreamHandler(g, method)
		}
	}

	serverName := MockPrefix + service.GoName + ServerSuffix

	// Server structure.
	tm.generateStruct(g, serverName)

	// NewServer factory.
	tm.generateNewFunc(g, service, serverName)

	// Server method implementations.
	for _, method := range service.Methods {
		tm.generateMethodDefinitions(g, tm.serverMethod(g, method))
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			tm.generateServerStreamHandler(g, method)
		}
	}
}

func (tm *testifyMocker) generateStruct(g *protogen.GeneratedFile, typeName string) {
	g.P("type ", unexport(typeName), " struct {")
	g.P(g.QualifiedGoIdent(testifyMockPackage.Ident("Mock")))
	g.P("}")
	g.P()
}

func (tm *testifyMocker) generateNewFunc(g *protogen.GeneratedFile, service *protogen.Service, typeName string) {
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P(deprecationComment)
	}
	g.P("func New", typeName, " (", ") *", unexport(typeName), " {")
	g.P("return &", unexport(typeName), "{}")
	g.P("}")
	g.P()
}

func (tm *testifyMocker) generateClientStreamHandler(g *protogen.GeneratedFile, method *protogen.Method) {
	clientStreamHandler := MockPrefix + method.Parent.GoName + "_" + method.GoName + ClientSuffix
	tm.generateStruct(g, clientStreamHandler)

	tm.generateNewFunc(g, method.Parent, clientStreamHandler)

	g.P("func (x *", unexport(clientStreamHandler), ") Header() (", g.QualifiedGoIdent(grpcMetaPackage.Ident("MD")), ", error) {")
	g.P("args := x.Called()")
	g.P("return args.Get(0).(", g.QualifiedGoIdent(grpcMetaPackage.Ident("MD")), "), args.Error(1)")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(clientStreamHandler), ") Trailer() ", g.QualifiedGoIdent(grpcMetaPackage.Ident("MD")), " {")
	g.P("args := x.Called()")
	g.P("return args.Get(0).(", g.QualifiedGoIdent(grpcMetaPackage.Ident("MD")), ")")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(clientStreamHandler), ") CloseSend() error {")
	g.P("args := x.Called()")
	g.P("return args.Error(0)")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(clientStreamHandler), ") Context() ", g.QualifiedGoIdent(contextPackage.Ident("Context")), " {")
	g.P("args := x.Called()")
	g.P("return args.Get(0).(", g.QualifiedGoIdent(contextPackage.Ident("Context")), ")")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(clientStreamHandler), ") SendMsg(m interface{}) error {")
	g.P("args := x.Called(m)")
	g.P("return args.Error(0)")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(clientStreamHandler), ") RecvMsg(m interface{}) error {")
	g.P("args := x.Called(m)")
	g.P("return args.Error(0)")
	g.P("}")
	g.P()

	if method.Desc.IsStreamingClient() {
		g.P("func (x *", unexport(clientStreamHandler), ") Send(m *", g.QualifiedGoIdent(method.Input.GoIdent), ") error {")
		g.P("args := x.Called(m)")
		g.P("return args.Error(0)")
		g.P("}")
		g.P()

		g.P("func (x *", unexport(clientStreamHandler), ") OnSend(m interface{}) *", g.QualifiedGoIdent(testifyMockPackage.Ident("Call")), " {")
		g.P("return x.On(\"Send\", m)")
		g.P("}")
		g.P()
	}

	methodName := "Recv"
	if !method.Desc.IsStreamingServer() {
		methodName = "CloseAndRecv"
	}

	g.P("func (x *", unexport(clientStreamHandler), ") ", methodName, "() ", "(*", g.QualifiedGoIdent(method.Output.GoIdent), ", error) {")
	g.P("args := x.Called()")
	g.P("return args.Get(0).(*", g.QualifiedGoIdent(method.Output.GoIdent), "), args.Error(1)")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(clientStreamHandler), ") On", methodName, "() *", g.QualifiedGoIdent(testifyMockPackage.Ident("Call")), " {")
	g.P("return x.On(\"", methodName, "\")")
	g.P("}")
	g.P()
}

func (tm *testifyMocker) generateServerStreamHandler(g *protogen.GeneratedFile, method *protogen.Method) {
	serverStreamHandler := MockPrefix + method.Parent.GoName + "_" + method.GoName + ServerSuffix
	tm.generateStruct(g, serverStreamHandler)

	tm.generateNewFunc(g, method.Parent, serverStreamHandler)

	g.P("func (x *", unexport(serverStreamHandler), ") SetHeader(md ", g.QualifiedGoIdent(grpcMetaPackage.Ident("MD")), ") error {")
	g.P("args := x.Called(md)")
	g.P("return args.Error(0)")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(serverStreamHandler), ") SendHeader(md ", g.QualifiedGoIdent(grpcMetaPackage.Ident("MD")), ") error {")
	g.P("args := x.Called(md)")
	g.P("return args.Error(0)")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(serverStreamHandler), ") SetTrailer(md ", g.QualifiedGoIdent(grpcMetaPackage.Ident("MD")), ") {")
	g.P("_ = x.Called(md)")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(serverStreamHandler), ") Context() ", g.QualifiedGoIdent(contextPackage.Ident("Context")), " {")
	g.P("args := x.Called()")
	g.P("return args.Get(0).(", g.QualifiedGoIdent(contextPackage.Ident("Context")), ")")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(serverStreamHandler), ") SendMsg(m interface{}) error {")
	g.P("args := x.Called(m)")
	g.P("return args.Error(0)")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(serverStreamHandler), ") RecvMsg(m interface{}) error {")
	g.P("args := x.Called(m)")
	g.P("return args.Error(0)")
	g.P("}")
	g.P()

	if method.Desc.IsStreamingClient() {
		g.P("func (x *", unexport(serverStreamHandler), ") Recv() ", "(*", g.QualifiedGoIdent(method.Input.GoIdent), ", error) {")
		g.P("args := x.Called()")
		g.P("return args.Get(0).(*", g.QualifiedGoIdent(method.Input.GoIdent), "), args.Error(1)")
		g.P("}")
		g.P()

		g.P("func (x *", unexport(serverStreamHandler), ") OnRecv() *", g.QualifiedGoIdent(testifyMockPackage.Ident("Call")), " {")
		g.P("return x.On(\"Recv\")")
		g.P("}")
		g.P()
	}

	methodName := "Send"
	if !method.Desc.IsStreamingServer() {
		methodName = "SendAndClose"
	}

	g.P("func (x *", unexport(serverStreamHandler), ") ", methodName, "(m *", g.QualifiedGoIdent(method.Output.GoIdent), ") error {")
	g.P("args := x.Called(m)")
	g.P("return args.Error(0)")
	g.P("}")
	g.P()

	g.P("func (x *", unexport(serverStreamHandler), ") On", methodName, "(m interface{}) *", g.QualifiedGoIdent(testifyMockPackage.Ident("Call")), " {")
	g.P("return x.On(\"", methodName, "\", m)")
	g.P("}")
	g.P()
}

func (tm *testifyMocker) generateMethodDefinitions(g *protogen.GeneratedFile, method *model.Method) {
	if method.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
		g.P(deprecationComment)
	}
	g.P(method, "{")
	args := make([]string, len(method.Arguments))
	for i, a := range method.Arguments {
		args[i] = a.Name
	}

	lastArg := *method.Arguments[len(method.Arguments)-1]
	if lastArg.Type.IsVariadic() {
		g.P(lastArg.Name, "0 := []interface{}{", strings.Join(args[:len(args)-1], ", "), "}")
		g.P("for _, ", lastArg.Name, "1 := range ", lastArg.Name, " {")
		g.P(lastArg.Name, "0 = append(", lastArg.Name, "0, ", lastArg.Name, "1)")
		g.P("}")
		g.P("args := ", method.Receiver.Name, ".Called(", lastArg.Name, "0...)")
	} else {
		g.P("args := ", method.Receiver.Name, ".Called(", strings.Join(args, ", "), ")")
	}

	if len(method.Return) > 0 {
		ret := make([]string, len(method.Return))
		for i, r := range method.Return {
			switch r {
			case "bool":
				ret[i] = fmt.Sprintf("args.Bool(%d)", i)
			case "int":
				ret[i] = fmt.Sprintf("args.Int(%d)", i)
			case "string":
				ret[i] = fmt.Sprintf("args.String(%d)", i)
			case "error":
				ret[i] = fmt.Sprintf("args.Error(%d)", i)
			default:
				ret[i] = fmt.Sprintf("args.Get(%d).(%v)", i, r)
			}
		}
		g.P("return ", strings.Join(ret, ", "))
	}

	g.P("}")
	g.P()

	methodName := method.GoName
	method.GoName = "On" + method.GoName
	method.SetReturn("*" + g.QualifiedGoIdent(testifyMockPackage.Ident("Call")))

	for _, arg := range method.Arguments {
		argType := "interface{}"
		if arg.Type.IsVariadic() {
			argType = "..." + argType
		}
		arg.Type = model.QualifiedGoIdent(argType)
	}

	g.P(method, "{")
	if lastArg.Type.IsVariadic() {
		g.P("return ", method.Receiver.Name, ".On(\"", methodName, "\", append([]interface{}{", strings.Join(args[:len(args)-1], ", "), "},", lastArg.Name, "...)...)")
	} else {
		g.P("return ", method.Receiver.Name, ".On(\"", methodName, "\", ", strings.Join(args, ", "), ")")
	}
	g.P("}")
	g.P()
}

func (tm *testifyMocker) clientMethod(g *protogen.GeneratedFile, method *protogen.Method) *model.Method {
	m := model.NewMethod(method, model.Receiver{Name: "c", Type: "*" + unexport(MockPrefix+method.Parent.GoName+ClientSuffix)})
	m.AddArgument("ctx", g.QualifiedGoIdent(contextPackage.Ident("Context")))
	if !method.Desc.IsStreamingClient() {
		m.AddArgument("in", "*"+g.QualifiedGoIdent(method.Input.GoIdent))
	}
	m.AddArgument("opts", "..."+g.QualifiedGoIdent(grpcPackage.Ident("CallOption")))
	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		m.AddReturn("*" + g.QualifiedGoIdent(method.Output.GoIdent))
	} else {
		m.AddReturn(g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       method.Parent.GoName + "_" + method.GoName + ClientSuffix,
			GoImportPath: method.Output.GoIdent.GoImportPath,
		}))
	}
	m.AddReturn("error")
	return m
}

func (tm *testifyMocker) serverMethod(g *protogen.GeneratedFile, method *protogen.Method) *model.Method {
	m := model.NewMethod(method, model.Receiver{Name: "s", Type: "*" + unexport(MockPrefix+method.Parent.GoName+ServerSuffix)})
	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		m.AddArgument("ctx", g.QualifiedGoIdent(contextPackage.Ident("Context")))
		m.AddReturn("*" + g.QualifiedGoIdent(method.Output.GoIdent))
	}
	if !method.Desc.IsStreamingClient() {
		m.AddArgument("in", "*"+g.QualifiedGoIdent(method.Input.GoIdent))
	}
	if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
		m.AddArgument("out", g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       method.Parent.GoName + "_" + method.GoName + ServerSuffix,
			GoImportPath: method.Output.GoIdent.GoImportPath,
		}))
	}
	m.AddReturn("error")
	return m
}

func init() {
	setMocker("testify", NewTestifyMocker)
}
