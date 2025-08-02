package main

import (
	"cloud-server/internal/handlers"
	"cloud-server/internal/models"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/upload", handler.UploadHandler).Methods("POST")
	api.HandleFunc("/delete", handler.DeleteHandler).Methods("DELETE")
	api.HandleFunc("/download", handler.FileDownloadHandler).Methods("GET")
	api.HandleFunc("/createDirectory", handler.CreateFolderHandler).Methods("POST")
	api.HandleFunc("/showTreeDirectory", handler.ListDirectoryTreeHandler).Methods("GET")
	api.HandleFunc("/createFolder", handler.CreateFolderHandler).Methods("POST")
	api.HandleFunc("/folder/{id:.*}", handler.DownloadFolderHandler).Methods("GET")

	storageRoot := filepath.Join(cfg.PATH_TO_DIRECTORY)
	api.PathPrefix("/files/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/files/")
		cleanPath := filepath.Clean(path)

		finalPath := filepath.Join(storageRoot, cleanPath)
		log.Println("Serving file:", finalPath)

		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(w, r, finalPath)
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
