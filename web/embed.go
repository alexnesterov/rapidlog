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
	sub, _ := fs.Sub(dist, "dist")
	return spaHandler(sub)
}

// spaHandler раздаёт файлы сборки, а на неизвестные пути без расширения (/login,
// /registration — маршруты клиентского роутера) отдаёт index.html. Отсутствующие файлы
// с расширением (/assets/old.js, /favicon.ico) остаются 404, а не превращаются в HTML.
func spaHandler(fsys fs.FS) http.Handler {
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
