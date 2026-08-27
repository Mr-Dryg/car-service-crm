package http

import "net/http"

func NewRouter(orderHandler *OrderHandler) *http.ServeMux {
	mainMux := http.NewServeMux()

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("POST /orders", orderHandler.Create)
	
	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	return mainMux
}
