package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/integronlabs/integron/array"
	"github.com/integronlabs/integron/helpers"
	httpOperation "github.com/integronlabs/integron/http"
	"github.com/integronlabs/integron/object"
	"github.com/integronlabs/integron/removenull"
	"github.com/integronlabs/integron/server"

	"github.com/swaggest/swgui/v5emb"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sirupsen/logrus"

	_ "embed"
)

//go:embed docs/openapi.yaml
var openapiSpec []byte

func main() {
	helpers.SetupLogging()

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromData(openapiSpec)
	if err != nil {
		panic(err)
	}

	// Validate document
	err = doc.Validate(ctx)
	if err != nil {
		panic(err)
	}

	r, err := gorillamux.NewRouter(doc)
	if err != nil {
		panic(err)
	}

	s := server.Server{
		Router:     r,
		LogHandler: func(ctx context.Context, r *http.Request) {},
	}

	server.RegisterStep("http", func(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
		client := http.Client{}
		return httpOperation.Run(ctx, &client, stepMap, stepOutputs)
	})
	server.RegisterStep("array", array.Run)
	server.RegisterStep("object", object.Run)
	server.RegisterStep("removenull", removenull.Run)
	server.RegisterStep("error", func(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
		return nil, "end", errors.New("error step triggered")
	})

	http.Handle("/", http.HandlerFunc(s.Handler))

	fs := http.FileServer(http.Dir("docs/"))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))

	http.Handle("/ui/", v5emb.New(
		"Integron Sunrise",
		"/docs/openapi.yaml",
		"/ui/",
	))

	http.Handle("/metrics", promhttp.Handler())

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
