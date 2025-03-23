package server

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/routers"
)

type Server struct {
	Router routers.Router
	Ctx    context.Context
	Client http.Client
}

type StepHandler func(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error)
