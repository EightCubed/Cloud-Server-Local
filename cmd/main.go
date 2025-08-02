package main

import (
	"cloud-server/internal/handlers"
	"cloud-server/internal/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

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

	cfg := &models.Config{
		PATH_TO_DIRECTORY: os.Getenv("FILE_STORAGE_PATH"),
		STORAGE_DIRECTORY: os.Getenv("FILE_STORAGE_DIRECTORY"),
	}

	r := mux.NewRouter()

	handler := handlers.ReturnHandler(logger, cfg)

	// Set up routes
	r.HandleFunc("/upload", handler.UploadHandler).Methods("POST")
	r.HandleFunc("/delete", handler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/download", handler.FileDownloadHandler).Methods("GET")
	r.HandleFunc("/createDirectory", handler.CreateFolderHandler).Methods("POST")
	r.HandleFunc("/showTreeDirectory", handler.ListDirectoryTreeHandler).Methods("GET")
	r.HandleFunc("/createFolder", handler.CreateFolderHandler).Methods("POST")
	r.HandleFunc("/folder/{id:.*}", handler.DownloadFolderHandler).Methods("GET")

	storageRoot := filepath.Join(cfg.PATH_TO_DIRECTORY)
	r.HandleFunc("/debug-path/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]

		log.Printf("Debug - Full URL: %s", r.URL.String())
		log.Printf("Debug - URL Path: %s", r.URL.Path)
		log.Printf("Debug - Mux Path Variable: %s", path)
		log.Printf("Debug - Storage Root: %s", storageRoot)

		finalPath := filepath.Join(storageRoot, path)
		log.Printf("Debug - Final Path: %s", finalPath)

		// Check if file exists
		if stat, err := os.Stat(finalPath); err == nil {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "File exists: %s\nSize: %d bytes\nModified: %s",
				finalPath, stat.Size(), stat.ModTime())
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "File not found: %s\nError: %s", finalPath, err)
		}
	}).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: c.Handler(r),
	}

	go func() {
		logger.Info("Server started at http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
