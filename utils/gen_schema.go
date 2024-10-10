package main

import (
	"os"
	"os/exec"

	"github.com/invopop/jsonschema"
	"github.com/sayden/counters"
	"github.com/stoewer/go-strcase"
)

func main() {
	r := new(jsonschema.Reflector)
	r.KeyNamer = strcase.SnakeCase

	schema := r.Reflect(&counters.CounterTemplate{})
	byt, err := schema.MarshalJSON()
	if err != nil {
		panic(err)
	}

	f, err := os.CreateTemp("/tmp", "schema.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// defer os.Remove(f.Name())

	_, err = f.Write(byt)
	if err != nil {
		panic(err)
	}

	exec.Command("/home/mcastro/.local/bin/generate-schema-doc", f.Name(), "docs/schema").Run()
}
