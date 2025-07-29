package main

import (
	"cloud-server/internal/handlers"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

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

	// fs := http.FileServer(http.Dir("./uploads"))

	// Set up routes
	r.HandleFunc("/upload", handlers.UploadHandler(logger)).Methods("POST")
	r.HandleFunc("/delete", handlers.DeleteHandler(logger)).Methods("DELETE")
	r.HandleFunc("/download", handlers.FileDownloadHandler(logger)).Methods("GET")
	r.HandleFunc("/createDirectory", handlers.CreateFolderHandler(logger)).Methods("POST")
	r.HandleFunc("/showTreeDirectory", handlers.ListFileDirectoryHandler(logger)).Methods("GET")
	r.HandleFunc("/listFiles", handlers.ListFilesByPath(logger)).Methods("GET")
	r.HandleFunc("/createFolder", handlers.CreateFolderHandler(logger)).Methods("POST")

	r.PathPrefix("/uploads/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/uploads/")
		filePath := filepath.Join("/home/rockon/Github/Cloud-Server-Local/uploads", path)

		mimeType := mime.TypeByExtension(filepath.Ext(filePath))
		if mimeType != "" {
			w.Header().Set("Content-Type", mimeType)
		}

		http.ServeFile(w, r, filePath)
	})

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
