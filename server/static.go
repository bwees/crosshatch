package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"crosshatch/internal/config"

	"github.com/go-fuego/fuego"
)

func registerStaticWeb(server *fuego.Server) {
	root := config.WebStaticPath()
	if root == "" {
		return
	}

	index := filepath.Join(root, "index.html")
	fileServer := http.FileServer(http.Dir(root))

	handler := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		clean := filepath.Join(root, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(clean); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, index)
	}

	// Register method-less ("/") rather than "GET /": Go's ServeMux treats a
	// method-specific catch-all as conflicting with the method-less routes
	// fuego registers (e.g. "/api/ws/"), so GetStd would panic at startup.
	fuego.Handle(server, "/", http.HandlerFunc(handler), fuego.OptionHide())
	slog.Info("Serving web frontend", "path", root)
}
