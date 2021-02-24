package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/ardifirmansyah/genshin-gacha-viewer/src/gatcha"
)

func main() {
	_ = http.ListenAndServe(":8080", newRouter())
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	fileServer(r, "/static", http.Dir("assets"))

	r.Get("/", gatcha.Index)
	r.Post("/gacha/process/history", gatcha.Process)

	return r
}

func fileServer(r *chi.Mux, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}