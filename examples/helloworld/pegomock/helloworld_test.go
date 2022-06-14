package helloworld

import (
	"context"
	reflect "reflect"
	"testing"

	"github.com/petergtz/pegomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func AnyContextContext() context.Context {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(context.Context))(nil)).Elem()))
	var nullValue context.Context
	return nullValue
}

func AnyGrpcCallOption() grpc.CallOption {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(grpc.CallOption))(nil)).Elem()))
	var nullValue grpc.CallOption
	return nullValue
}

func TestSayHello(t *testing.T) {
	// Create a new mock client for the Greeter service.
	m := NewMockGreeterClient()

	// Create the request and response.
	ctx := context.Background()
	req := &HelloRequest{Name: "Felix"}
	res := &HelloReply{Message: "Hello, world!"}

	// Set up the expectation.
	pegomock.When(m.SayHello(ctx, req)).ThenReturn(res, nil)

	// Call the client.
	r, err := m.SayHello(ctx, req)

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)
}

func TestSayHelloWithOptions(t *testing.T) {
	// Create a new mock client for the Greeter service.
	m := NewMockGreeterClient()

	// Create the request and response.
	ctx := context.Background()
	req := &HelloRequest{Name: "Felix"}
	res := &HelloReply{Message: "Hello, world!"}

	// Set up the expectation.
	pegomock.When(m.SayHello(ctx, req, grpc.WaitForReady(true))).ThenReturn(res, nil)

	// Call the client.
	r, err := m.SayHello(ctx, req, grpc.WaitForReady(true))

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)
}

func TestSayHelloWithAnyOptions(t *testing.T) {
	// Create a new mock client for the Greeter service.
	m := NewMockGreeterClient()

	// Create the request and response.
	ctx := context.Background()
	req := &HelloRequest{Name: "Felix"}
	res := &HelloReply{Message: "Hello, world!"}

	// Set up the expectation.
	pegomock.When(m.SayHello(AnyContextContext(), AnyPtrToHelloworldHelloRequest(), AnyGrpcCallOption())).ThenReturn(res, nil)

	// Call the client.
	r, err := m.SayHello(ctx, req, grpc.WaitForReady(true))

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)
}
