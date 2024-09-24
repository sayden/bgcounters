package main

import (
	"github.com/alecthomas/kong"
	"github.com/sayden/counters/pipelines"
)

type vassal struct {
	pipelines.VassalConfig
}

func (i *vassal) Run(ctx *kong.Context) error {
	return pipelines.CSVToVassalFile(i.VassalConfig)
}
