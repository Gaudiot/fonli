package health

import (
	"encoding/json"
	"net/http"
	"time"
)

type response struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Handler responde a GET com JSON indicando que o processo está vivo (liveness).
// Não executa verificações de dependências (BD, etc.); use um endpoint /ready separado se precisar disso.
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	}
}
