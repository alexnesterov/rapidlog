// Package web exposes the production frontend bundled into the API binary.
package web

import "embed"

//go:embed dist
var Dist embed.FS
