package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pintoter/persons/internal/service"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	router  *mux.Router
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	handler := &Handler{
		router:  mux.NewRouter(),
		service: service,
	}

	handler.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
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
		v1.HandleFunc("/persons", h.createPerson).Methods(http.MethodPost)
		v1.HandleFunc("/persons/{id:[0-9]+}", h.getPerson).Methods(http.MethodGet)
		v1.HandleFunc("/persons", h.getPersons).Methods(http.MethodGet)
		v1.HandleFunc("/persons/{id:[0-9]+}", h.updatePerson).Methods(http.MethodPatch)
		v1.HandleFunc("/persons/{id:[0-9]+}", h.deletePerson).Methods(http.MethodDelete)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
