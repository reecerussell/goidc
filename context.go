package goidc

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

// Context is an implementation of context.Context, used to
// pass application variables throughout various processes.
type Context struct {
	ctx context.Context
}

// NewContext returns a new instance of *Context, with values from req.
func NewContext(ctx context.Context, req *events.APIGatewayProxyRequest) context.Context {
	for key, value := range req.StageVariables {
		ck := NewContextKey(fmt.Sprintf("STAGE:%s", key))
		ctx = context.WithValue(ctx, ck, value)
	}

	return &Context{
		ctx: ctx,
	}
}

// Deadline is a wrapper around the underlying context, to implement the interface.
func (ctx *Context) Deadline() (time.Time, bool) {
	return ctx.ctx.Deadline()
}

// Done is a wrapper around the underlying context, to implement the interface.
func (ctx *Context) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

// Err is a wrapper around the underlying context, to implement the interface.
func (ctx *Context) Err() error {
	return ctx.ctx.Err()
}

// Value is a wrapper around the underlying context, to implement the interface.
func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.ctx.Value(key)
}

// ContextKey is an abstraction of string, used to satisfy the need of a complex
// type when accessing a context's values.
type ContextKey string

// NewContextKey returns a new instance of ContextKey for the given key.
func NewContextKey(key string) ContextKey {
	return ContextKey(key)
}

// StageVariable returns a stage variable from the Context.
// This will panic if the variable does not exist.
func StageVariable(ctx context.Context, key string) string {
	ck := NewContextKey(fmt.Sprintf("STAGE:%s", key))
	value := ctx.Value(ck)
	if value == nil {
		panic(fmt.Errorf("the stage variable '%s' was not found", key))
	}

	return value.(string)
}
