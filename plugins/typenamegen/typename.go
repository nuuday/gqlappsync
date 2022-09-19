package typenamegen

import (
	_ "embed"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
)

//go:embed typename.gotpl
var typenameTemplate string

func New() plugin.Plugin {
	return &Plugin{}
}

type Plugin struct {
}

var (
	_ plugin.CodeGenerator = &Plugin{}
	_ plugin.ConfigMutator = &Plugin{}
)

func (m *Plugin) Name() string {
	return "typenamegen"
}

func (m *Plugin) MutateConfig(cfg *config.Config) error {

	_ = syscall.Unlink(filename(cfg.Model))
	return nil
}

type TypenameData struct {
	Objects                                             codegen.Objects
	Interfaces                                          map[string]*codegen.Interface
	TypenamesDirectlyOrIndirectlyConnectedToAnInterface map[string]bool
	TypePrefixes                                        []string
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {

	typesDirectlyOrIndirectlyConnectedToAnInterface := createMapOfTypenamesDirectlyOrIndirectlyConnectedToAnInterface(data)

	t := TypenameData{
		Objects:    data.Objects,
		Interfaces: data.Interfaces,
		TypenamesDirectlyOrIndirectlyConnectedToAnInterface: typesDirectlyOrIndirectlyConnectedToAnInterface,
		TypePrefixes: []string{"", "*"},
	}
	return templates.Render(templates.Options{
		PackageName:     data.Config.Model.Package,
		Filename:        filename(data.Config.Model),
		Data:            t,
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
		Template:        typenameTemplate,
	})
}

func filename(pkgCfg config.PackageConfig) string {
	dir := filepath.Dir(pkgCfg.Filename)
	filename := filepath.Join(dir, "typenames_gen.go")
	return filename
}

func createMapOfTypenamesDirectlyOrIndirectlyConnectedToAnInterface(data *codegen.Data) map[string]bool {

	parentsByChildTypename := make(map[string][]string)
	for _, object := range data.Objects {
		objectTypename := object.Definition.Name

		if objectTypename == "Query" || objectTypename == "Mutation" || objectTypename == "Subscription" {
			continue
		}
		for _, field := range object.Fields {
			if field.TypeReference.Definition.Kind == "SCALAR" {
				continue
			}
			if field.TypeReference.Definition.Kind == "ENUM" {
				continue
			}
			fieldTypename := field.TypeReference.Definition.Name
			if strings.HasPrefix(fieldTypename, "__") {
				continue
			}
			val, keyExists := parentsByChildTypename[fieldTypename]
			if !keyExists {
				parentsByChildTypename[fieldTypename] = []string{}
			}

			parentsByChildTypename[fieldTypename] = append(val, objectTypename)
		}
	}
	typesDirectlyOrIndirectlyConnectedToAnInterface := make(map[string]bool)
	for _, e := range data.Interfaces {
		interfaceName := e.Name
		queue := []string{interfaceName}
		for len(queue) > 0 {
			item, rest := queue[0], queue[1:]
			slice, ok := parentsByChildTypename[item]
			if ok {
				for _, ee := range slice {
					_, keyExists := typesDirectlyOrIndirectlyConnectedToAnInterface[ee]
					if !keyExists {
						typesDirectlyOrIndirectlyConnectedToAnInterface[ee] = true
					}
					rest = append(rest, ee)
				}
			}
			queue = rest
		}
	}
	return typesDirectlyOrIndirectlyConnectedToAnInterface
}
