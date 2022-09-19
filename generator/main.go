package generator

import (
	"errors"
	"go/types"
	"reflect"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/nuuday/gqlappsync/plugins/typenamegen"
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

func Run(configFilename string) error {
	var cfg *config.Config
	var err error
	if configFilename != "" {
		cfg, err = config.LoadConfig(configFilename)
		if err != nil {
			return err
		}
	} else {
		if cfg, err = config.LoadConfigFromDefaultLocations(); err != nil {
			return errors.New("couldn't load gqlgen.yaml from current location")
		}
	}
	if err := generate(cfg); err != nil {
		panic(err)
	}
	return nil
}
func generate(cfg *config.Config) error {
	plugins := []plugin.Plugin{
		&modelgen.Plugin{
			MutateHook: mutateHook,
		},
		typenamegen.New(),
	}
	if len(cfg.SchemaFilename) == 0 {
		return errors.New("schema file not found")
	}
	if err := cfg.LoadSchema(); err != nil {
		return err
	}

	if err := cfg.Init(); err != nil {
		return err
	}

	for _, p := range plugins {
		if mut, ok := p.(plugin.ConfigMutator); ok {
			err := mut.MutateConfig(cfg)
			if err != nil {
				return err
			}
		}
	}

	data, err := codegen.BuildData(cfg)
	if err != nil {
		return err
	}

	for _, p := range plugins {
		if mut, ok := p.(plugin.CodeGenerator); ok {
			err := mut.GenerateCode(data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
