package main

import (
	"cloud-server/internal/handlers"
	"cloud-server/internal/models"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := &models.Config{
		PATH_TO_DIRECTORY: os.Getenv("FILE_STORAGE_PATH"),
		STORAGE_DIRECTORY: os.Getenv("FILE_STORAGE_DIRECTORY"),
	}

	r := mux.NewRouter()

	handler := handlers.ReturnHandler(logger, cfg)

	// Set up routes
	r.HandleFunc("/upload", handler.UploadHandler).Methods("POST")
	// r.HandleFunc("/delete", handlers.DeleteHandler(logger)).Methods("DELETE")
	// r.HandleFunc("/download", handlers.FileDownloadHandler(logger)).Methods("GET")
	r.HandleFunc("/createDirectory", handler.CreateFolderHandler).Methods("POST")
	r.HandleFunc("/showTreeDirectory", handler.ListDirectoryTreeHandler).Methods("GET")
	// r.HandleFunc("/listFiles", handlers.ListFilesByPath(logger)).Methods("GET")
	// r.HandleFunc("/createFolder", handlers.CreateFolderHandler(logger)).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	// Start the HTTP server on port 8080
	logger.Info("Server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", c.Handler(r))
	if err != nil {
		logger.Fatal("Error starting server:", zap.Error(err))
	}
}
