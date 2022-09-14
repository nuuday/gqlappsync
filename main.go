package main

import (
	"errors"
	"go/types"
	"io/fs"
	"reflect"

	"github.com/nuuday/gqlappsync/plugins/typenamegen"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"

	"github.com/99designs/gqlgen/plugin/modelgen"
)

func mutateHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {
		metaValue := reflect.ValueOf(model).Elem()
		if metaValue.FieldByName("Implements").Len() > 0 { // adds Typename if type implements an interface
			model.Fields = append(model.Fields, &modelgen.Field{
				GoName: "Typename",
				Type:   types.Typ[types.String].Underlying(),
				Tag:    `json:"__typename"`,
			})
		}
	}

	return b
}

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.LoadConfigFromDefaultLocations()
	if errors.Is(err, fs.ErrNotExist) {
		cfg, err = config.LoadDefaultConfig()
	}
	if err != nil {
		panic(err)
	}

	plugins := []plugin.Plugin{
		&modelgen.Plugin{
			MutateHook: mutateHook,
		},
		typenamegen.New(),
	}

	if err := cfg.LoadSchema(); err != nil {
		panic(err)
	}

	if err := cfg.Init(); err != nil {
		panic(err)
	}

	for _, p := range plugins {
		if mut, ok := p.(plugin.ConfigMutator); ok {
			err := mut.MutateConfig(cfg)
			if err != nil {
				panic(err)
			}
		}
	}

	data, err := codegen.BuildData(cfg)
	if err != nil {
		panic(err)
	}

	for _, p := range plugins {
		if mut, ok := p.(plugin.CodeGenerator); ok {
			err := mut.GenerateCode(data)
			if err != nil {
				panic(err)
			}
		}
	}
}
