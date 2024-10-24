package main

import (
	"cloud-server/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	r := mux.NewRouter()

	// Set up routes
	r.HandleFunc("/upload", handlers.UploadHandler(logger)).Methods("POST")
	r.HandleFunc("/download", handlers.FileDownloadHandler(logger)).Methods("GET")
	r.HandleFunc("/createDirectory", handlers.CreateFolderHandler(logger)).Methods("POST")
	r.HandleFunc("/showTreeDirectory", handlers.ListFileDirectoryHandler(logger)).Methods("GET")
	r.HandleFunc("/listFiles", handlers.ListFilesByPath(logger)).Methods("GET")
	r.HandleFunc("/createFolder", handlers.CreateFolderHandler(logger)).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Specify the allowed origin
		AllowCredentials: false,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	// Start the HTTP server on port 8080
	logger.Info("Server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", c.Handler(r))
	if err != nil {
		logger.Fatal("Error starting server:", zap.Error(err))
	}
}
