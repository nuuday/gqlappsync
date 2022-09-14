package typenamegen

import (
	_ "embed"
	"path/filepath"
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

func (m *Plugin) GenerateCode(data *codegen.Data) error {

	return templates.Render(templates.Options{
		PackageName:     data.Config.Model.Package,
		Filename:        filename(data.Config.Model),
		Data:            data,
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
