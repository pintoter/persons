package transport

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pintoter/persons/internal/service"
	"github.com/pintoter/persons/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Config interface {
	GetAddr() string
}

type Handler struct {
	router  *mux.Router
	service *service.Service
}

func NewHandler(service *service.Service, cfg Config) *Handler {
	handler := &Handler{
		router:  mux.NewRouter(),
		service: service,
	}

	swaggerURL := fmt.Sprintf("http://%s/swagger/doc.json", cfg.GetAddr())

	handler.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerURL),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	handler.InitRoutes()

	return handler
}

func (h *Handler) InitRoutes() {
	v1 := h.router.PathPrefix("/api/v1").Subrouter()
	{
		v1.HandleFunc("/person", h.createNote).Methods(http.MethodPost)
		v1.HandleFunc("/person/{id:[0-9]+}", h.getNote).Methods(http.MethodGet)
		v1.HandleFunc("/person/{id:[0-9]+}", h.updateNote).Methods(http.MethodPatch)
		v1.HandleFunc("/person/{id:[0-9]+}", h.deleteNote).Methods(http.MethodDelete)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.DebugKV(r.Context(), "new request", "Method", r.Method, "url", r.URL, "addr:", r.RemoteAddr)
	h.router.ServeHTTP(w, r)
}
