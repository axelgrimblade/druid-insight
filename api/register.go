package api

import (
	"druid-insight/auth"
	"druid-insight/config"
	"druid-insight/logging"
	"net/http"
	"os"
)

func RegisterHandlers(cfg *auth.Config, users *auth.UsersFile, druidCfg *config.DruidConfig, accessLogger, loginLogger, reportLogger *logging.Logger) {
	http.HandleFunc("/api/login", withCORS(LoginHandler(cfg, users, loginLogger)))
	http.HandleFunc("/api/schema", withCORS(SchemaHandler(cfg, druidCfg, accessLogger)))
	http.HandleFunc("/api/reports/execute", withCORS(ReportExecuteHandler(cfg, users, druidCfg, accessLogger)))
	http.HandleFunc("/api/reports/status", withCORS(ReportStatusHandler(cfg)))
	http.HandleFunc("/api/reports/download", withCORS(DownloadReportCSV(cfg)))
	http.HandleFunc("/api/filters/values", withCORS(GetDimensionValues(cfg, druidCfg)))

	// Register /api/doc endpoint only if configured and file exists
	if cfg.DocFile != "" {
		if stat, err := os.Stat(cfg.DocFile); err == nil && !stat.IsDir() {
			http.HandleFunc("/api/doc", withCORS(DocDownloadHandler(cfg, cfg.DocFile)))
		}
	}
}

func StartServer(listenAddr string) error {
	return http.ListenAndServe(listenAddr, nil)
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}
