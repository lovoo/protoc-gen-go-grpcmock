package model

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// A QualifiedGoIdent is the qualified Go identifier of a type.
// For example: `context.Context` or `string`.
//
// If the identifier is from a different Go package it will be
// qualified (package.name) and an import statement for the
// identifier's package is expected be included in the file.
// The identifier may also be a builtin or pointer type.
type QualifiedGoIdent string

// IsPointer indicated if the go identifier is a pointer.
func (ident QualifiedGoIdent) IsPointer() bool { return strings.HasPrefix(string(ident), "*") }

// IsVariadic indicates if the go identifier is variadic.
func (ident QualifiedGoIdent) IsVariadic() bool { return strings.HasPrefix(string(ident), "...") }

// Argument is the identifier of a function argument, consisting of a name and type.
// The type is the qualified Go identifier.
// For example: `ctx context.Context` => Argument{Name: "ctx", Type: "context.Context"}.
type Argument struct {
	Name string
	Type QualifiedGoIdent
}

// String implements the fmt.Stringer interface for the Argument type.
// It simply combines the name and type with a space in between, for example: `ctx context.Context`.
func (arg Argument) String() string { return fmt.Sprintf("%s %s", arg.Name, arg.Type) }

// Receiver is the identifier of a type, a method is defined on. It consists of the
// name and type of the receiver.
type Receiver struct {
	Name string
	Type string
}

// String implements the fmt.Stringer interface for the Receiver type.
// It simply combines the name and type with a space in between, for example: `c *MyClient`.
func (rec Receiver) String() string { return fmt.Sprintf("%s %s", rec.Name, rec.Type) }

// A Method describes a method in a service.
// It extends the protogen's Method with the additional fields.
type Method struct {
	*protogen.Method

	Receiver  Receiver
	Arguments []*Argument
	Return    []QualifiedGoIdent
}

// NewMethod creates a new Method based on the provided protogen.Method.
func NewMethod(method *protogen.Method, receiver Receiver) *Method {
	protoMethod := *method
	m := &Method{
		Method:    &protoMethod,
		Receiver:  receiver,
		Arguments: make([]*Argument, 0),
		Return:    make([]QualifiedGoIdent, 0),
	}
	return m
}

// AddArgument adds an argument to the method.
func (m *Method) AddArgument(name, ident string) *Method {
	m.Arguments = append(m.Arguments, &Argument{Name: name, Type: QualifiedGoIdent(ident)})
	return m
}

// AddReturn adds a return type to the method.
func (m *Method) AddReturn(ret string) *Method {
	m.Return = append(m.Return, QualifiedGoIdent(ret))
	return m
}

// SetReturn sets the return type of the method.
func (m *Method) SetReturn(ret string) *Method {
	m.Return = []QualifiedGoIdent{QualifiedGoIdent(ret)}
	return m
}

// String implements the fmt.Stringer interface for the method type.
// It formattes a valid Go function signature, for example:
// `func (c *MyClient) MyMethod(ctx context.Context, in *MyRequest, opts ...grpc.CallOption) (*MyResponse, error)`.
func (m *Method) String() string {
	args := make([]string, len(m.Arguments))
	for i, arg := range m.Arguments {
		args[i] = arg.String()
	}

	switch len(m.Return) {
	case 0:
		return fmt.Sprintf("func (%s) %s(%s)", m.Receiver, m.GoName, strings.Join(args, ", "))
	case 1:
		return fmt.Sprintf("func (%s) %s(%s) %v", m.Receiver, m.GoName, strings.Join(args, ", "), m.Return[0])
	default:
		ret := make([]string, len(m.Return))
		for i, r := range m.Return {
			ret[i] = string(r)
		}
		return fmt.Sprintf("func (%s) %s(%s) (%s)", m.Receiver, m.GoName, strings.Join(args, ", "), strings.Join(ret, ", "))
	}
}
