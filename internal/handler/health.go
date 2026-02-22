package handler

import (
	"encoding/json"
	"net/http"

	"github.com/YOUR_ORG/myapp/internal/store"
)

func health(st store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"status":   "ok",
			"database": "disabled",
		}

		if st != nil {
			if err := st.Ping(r.Context()); err != nil {
				resp["status"] = "degraded"
				resp["database"] = "error"
			} else {
				resp["database"] = "ok"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if resp["status"] != "ok" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}
}
