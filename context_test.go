package goidc

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestContextValue_GivenStageVariableKey_ReturnsValue(t *testing.T) {
	req := &events.APIGatewayProxyRequest{
		StageVariables: map[string]string{
			"env": "test",
		},
	}

	ctx := context.Background()
	ctx = NewContext(ctx, req)

	value := ctx.Value(NewContextKey("STAGE:env"))
	assert.Equal(t, "test", value)
}

func TestContextDeadline_GivenCtxWithDeadline_ReturnsDeadline(t *testing.T) {
	req := &events.APIGatewayProxyRequest{
		StageVariables: map[string]string{
			"env": "test",
		},
	}

	testDeadline := time.Now().Add(1 * time.Hour)
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, testDeadline)
	defer cancel()

	ctx = NewContext(ctx, req)

	deadline, ok := ctx.Deadline()
	assert.Equal(t, testDeadline, deadline)
	assert.True(t, ok)
}

func TestContextDone_GivenCtx_ReturnsChan(t *testing.T) {
	req := &events.APIGatewayProxyRequest{
		StageVariables: map[string]string{
			"env": "test",
		},
	}

	ctx := context.Background()
	testDone := ctx.Done()

	ctx = NewContext(ctx, req)
	done := ctx.Done()

	assert.Equal(t, testDone, done)
}

func TestContextErr_GivenCtxWithNoErr_ReturnsNil(t *testing.T) {
	req := &events.APIGatewayProxyRequest{
		StageVariables: map[string]string{
			"env": "test",
		},
	}

	ctx := context.Background()
	ctx = NewContext(ctx, req)

	err := ctx.Err()

	assert.Nil(t, err)
}

func TestStageVariable_GivenValidKey_ReturnsValue(t *testing.T) {
	req := &events.APIGatewayProxyRequest{
		StageVariables: map[string]string{
			"env": "test",
		},
	}

	ctx := context.Background()
	ctx = NewContext(ctx, req)

	assert.Equal(t, "test", StageVariable(ctx, "env"))
}

func TestStageVariable_GivenInvalidKey_Panics(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	req := &events.APIGatewayProxyRequest{
		StageVariables: map[string]string{},
	}

	ctx := context.Background()
	ctx = NewContext(ctx, req)

	_ = StageVariable(ctx, "env")
}
