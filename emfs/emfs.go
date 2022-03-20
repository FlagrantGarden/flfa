package emfs

import "embed"

//go:embed modules/***
var EmbeddedModulesFS embed.FS

func GetEmbeddedModulesFS() embed.FS {
	return EmbeddedModulesFS
}
