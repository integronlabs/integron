package main

import (
	"context"
	"errors"
	"flag"
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

	"github.com/sirupsen/logrus"

	_ "embed"
)

func main() {
	helpers.SetupLogging()

	openapiSpecPath := flag.String("spec", "docs/openapi.yaml", "Path to the OpenAPI spec")
	flag.Parse()

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(*openapiSpecPath)
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
		Router:       r,
		LogFormatter: &logrus.JSONFormatter{},
	}

	server.RegisterStep("http", func(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
		client := http.Client{}
		return httpOperation.Run(ctx, &client, stepMap, stepOutputs)
	})
	server.RegisterStep("transformarray", array.Run)
	server.RegisterStep("transformobject", object.Run)
	server.RegisterStep("removenull", removenull.Run)
	server.RegisterStep("error", func(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
		return nil, "end", errors.New("error step triggered")
	})

	http.Handle("/", http.HandlerFunc(s.Handler))

	fs := http.FileServer(http.Dir("docs/"))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))

	http.Handle("/ui/", v5emb.New(
		"Integron",
		"/docs/openapi.yaml",
		"/ui/",
	))

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
