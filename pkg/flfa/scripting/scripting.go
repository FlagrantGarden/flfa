package scripting

import "github.com/FlagrantGarden/flfa/pkg/tympan/module/scripting"

func NewEngine(modules []scripting.Module, libraries []scripting.Library) *scripting.Engine {
	engine := scripting.NewEngine()
	// ignore errors for now
	engine.SetStandardLibraries(engine.AllowedStandardLibraries())
	engine.AddApplicationLibraries(libraries...)
	for _, module := range modules {
		engine.AddApplicationModule(module)
	}
	return engine
}
