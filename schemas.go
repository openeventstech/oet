package main

import (
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

var (
	schEvent     *jsonschema.Schema
	schOrganizer *jsonschema.Schema
)

func loadSchemas() error {
	c := jsonschema.NewCompiler()
	var err error

	if schEvent, err = c.Compile("schemas/event.v1alpha1.json"); err != nil {
		return err
	}

	schOrganizer, err = c.Compile("schemas/organizer.v1alpha1.json")
	return err
}
