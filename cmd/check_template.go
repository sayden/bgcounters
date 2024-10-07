package main

import "github.com/alecthomas/kong"

type CheckTemplate struct {
	InputPath string `help:"Input path of the file to read" short:"i" required:"true"`
}

func (c *CheckTemplate) Run(ctx *kong.Context) error {
	return validateSchema(c.InputPath)
}
