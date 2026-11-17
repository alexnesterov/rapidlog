// Package web exposes the production frontend bundled into the API binary.
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed all:dist
var dist embed.FS

func Handler() http.Handler {
	fsys, _ := fs.Sub(dist, "dist")
	files := http.FileServerFS(fsys)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name != "" && path.Ext(name) == "" {
			if _, err := fs.Stat(fsys, name); err != nil {
				http.ServeFileFS(w, r, fsys, "index.html")
				return
			}
		}

		files.ServeHTTP(w, r)
	})
}
