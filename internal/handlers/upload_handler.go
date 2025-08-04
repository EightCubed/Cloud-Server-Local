package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (h Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger
	logger.Info("Calling upload handler")

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Unable to create multipart reader: "+err.Error(), http.StatusBadRequest)
		return
	}

	destPath := strings.TrimRight(r.URL.Query().Get("destPath"), "/")
	if destPath == "" {
		destPath = "Storage"
	}

	validatedPath, err := validatePathName(destPath, h.Config)
	if err != nil {
		http.Error(w, "Failed to validate pathname: "+err.Error(), http.StatusBadRequest)
		return
	}

	absolutePath := returnAbsolutePath(validatedPath, h.Config)
	if _, err := os.ReadDir(absolutePath); err != nil {
		err = os.Mkdir(absolutePath, 0750)
		if err != nil {
			http.Error(w, "Failed to create folder: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FormName() != "file" {
			continue
		}

		filename := part.FileName()
		if filename == "" {
			http.Error(w, "Filename not found in form part", http.StatusBadRequest)
			return
		}

		fullPath := filepath.Join(absolutePath, filename)
		out, err := os.Create(fullPath)
		if err != nil {
			http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(out, part) // streams directly to disk
		out.Close()
		if err != nil {
			http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully to %s", absolutePath)
}
