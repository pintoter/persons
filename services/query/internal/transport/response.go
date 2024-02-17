package transport

import (
	"encoding/json"
	"net/http"

	"github.com/pintoter/persons/pkg/logger"
	"github.com/pintoter/persons/services/query/internal/entity"
)

type getPersonResponse struct {
	Person entity.Person `json:"person"`
}

type getPersonsResponse struct {
	Persons []entity.Person `json:"persons"`
}

type errorResponse struct {
	Err string `json:"error"`
}

func renderJSON(w http.ResponseWriter, r *http.Request, code int, data any) {
	logger.DebugKV(r.Context(), "New response", "Code", code, "Response", data)
	resp, _ := json.MarshalIndent(data, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(resp)
}
