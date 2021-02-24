package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"

	"github.com/ardifirmansyah/genshin-gacha-viewer/src/gatcha"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = http.ListenAndServe(":" + port, newRouter())
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	fileServer(r, "/static", http.Dir("assets"))

	r.Get("/", gatcha.Index)
	r.Post("/gacha/process/history", gatcha.Process)

	r.Get("/ping", ping)

	return r
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
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