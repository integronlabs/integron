package server

import (
	"context"

	"github.com/getkin/kin-openapi/routers"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router       routers.Router
	LogFormatter logrus.Formatter
}

type StepHandler func(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error)
