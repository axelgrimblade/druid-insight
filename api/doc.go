package api

import (
	"druid-insight/auth"
	"net/http"
	"os"
	"path/filepath"
)

// DocDownloadHandler allows downloading a documentation file configured in config.yaml
func DocDownloadHandler(cfg *auth.Config, docFilePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if docFilePath == "" {
			http.Error(w, "Documentation not configured", http.StatusNotFound)
			return
		}
		absPath, err := filepath.Abs(docFilePath)
		if err != nil || absPath == "" {
			http.Error(w, "Invalid documentation path", http.StatusInternalServerError)
			return
		}
		if _, err := os.Stat(absPath); err != nil {
			http.Error(w, "Documentation file not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(absPath)+"\"")
		http.ServeFile(w, r, absPath)
	}
}
