package server

import (
	"common/logging"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewSwaggerServe() {
	logger := logging.WithComponent(context.Background(), "swagger")
	r := mux.NewRouter()

	// Route để phục vụ file swagger.json
	r.PathPrefix("/chat.swagger.json").Handler(http.FileServer(http.Dir("./proto/docs/chat")))

	// Tích hợp Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/chat.swagger.json"), // Trỏ đến file swagger.json
	))

	// API mẫu
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "pong"}`))
	}).Methods("GET")

	// Khởi chạy server
	logger.Info("Server is running at :3001")
	err := http.ListenAndServe(":3001", r)
	logger.Error(fmt.Sprintf("Swagger server stopped: %v", err))
	os.Exit(1)
}
