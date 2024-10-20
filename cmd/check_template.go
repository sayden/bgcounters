package main

import (
	"github.com/alecthomas/kong"
	"github.com/sayden/counters"
)

type CheckTemplate struct {
	InputPath      string `help:"Input path of the file to read" short:"i" required:"true"`
	IsCardTemplate bool   `help:"Check against the Card template schema instead of the counter schema" short:"c"`
}

func (c *CheckTemplate) Run(ctx *kong.Context) error {
	if c.IsCardTemplate {
		return counters.ValidateSchemaAtPath[counters.CardsTemplate](c.InputPath)
	}
	return counters.ValidateSchemaAtPath[counters.CounterTemplate](c.InputPath)
}
