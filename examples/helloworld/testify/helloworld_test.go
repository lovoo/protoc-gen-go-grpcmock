package helloworld

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestSayHello(t *testing.T) {
	// Create a new mock client for the Greeter service.
	m := NewMockGreeterClient()
	defer m.AssertExpectations(t)

	// Create the request and response.
	ctx := context.Background()
	req := &HelloRequest{Name: "Felix"}
	res := &HelloReply{Message: "Hello, world!"}

	// Set up the expectation.
	m.OnSayHello(ctx, req).Return(res, nil)

	// Call the client.
	r, err := m.SayHello(ctx, req)

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)
}

func TestSayHelloWithOptions(t *testing.T) {
	// Create a new mock client for the Greeter service.
	m := NewMockGreeterClient()
	defer m.AssertExpectations(t)

	// Create the request and response.
	ctx := context.Background()
	req := &HelloRequest{Name: "Felix"}
	res := &HelloReply{Message: "Hello, world!"}

	// Set up the expectation.
	m.OnSayHello(ctx, req, grpc.WaitForReady(true)).Return(res, nil)

	// Call the client.
	r, err := m.SayHello(ctx, req, grpc.WaitForReady(true))

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)
}

func TestSayHelloWithAnyOptions(t *testing.T) {
	// Create a new mock client for the Greeter service.
	m := NewMockGreeterClient()
	defer m.AssertExpectations(t)

	// Create the request and response.
	ctx := context.Background()
	req := &HelloRequest{Name: "Felix"}
	res := &HelloReply{Message: "Hello, world!"}

	// Set up the expectation.
	m.OnSayHello(ctx, AnyHelloRequest(), mock.Anything).Return(res, nil)

	// Call the client.
	r, err := m.SayHello(ctx, req, grpc.WaitForReady(true))

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)
}
