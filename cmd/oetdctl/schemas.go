package main

import (
	"bytes"
	"embed"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schemas
var schemas embed.FS

var (
	schEventV1A1     *jsonschema.Schema
	schLocationV1A1  *jsonschema.Schema
	schOrganizerV1A1 *jsonschema.Schema
)

func loadSchemas() error {
	var err error

	if schEventV1A1, err = loadSchema("event.v1alpha1.json"); err != nil {
		return err
	}
	if schLocationV1A1, err = loadSchema("location.v1alpha1.json"); err != nil {
		return err
	}
	schOrganizerV1A1, err = loadSchema("organizer.v1alpha1.json")
	return err
}

func loadSchema(name string) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()

	rawSchema, err := schemas.ReadFile("schemas/" + name)
	if err != nil {
		return nil, err
	}

	schema, err := jsonschema.UnmarshalJSON(bytes.NewReader(rawSchema))
	if err != nil {
		return nil, err
	}

	if err := compiler.AddResource(name, schema); err != nil {
		return nil, err
	}

	return compiler.Compile(name)
}
