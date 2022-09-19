package main

import (
	"fmt"
	"os"

	"github.com/nuuday/gqlappsync/generator"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "gqlappsync"
	app.Usage = "generate a graphql models based on schema"
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "config", Usage: "the config filename"},
	}
	app.Action = func(ctx *cli.Context) error {
		configFilename := ctx.String("config")
		if err := generator.Run(configFilename); err != nil {
			return err
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
