package server

import (
	"context"

	"github.com/getkin/kin-openapi/routers"
)

type Server struct {
	Router routers.Router
}

type StepHandler func(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error)
