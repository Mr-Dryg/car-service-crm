package http

import "net/http"

func NewRouter(orderHandler *OrderHandler) *http.ServeMux {
	mainMux := http.NewServeMux()

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("POST /orders", orderHandler.CreateFromUser)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("POST /orders", orderHandler.CreateFromManager)

	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))
	mainMux.Handle("/admin/", http.StripPrefix("/admin", adminMux))

	return mainMux
}
