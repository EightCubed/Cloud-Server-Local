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

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	destPath := r.FormValue("destPath")
	if destPath == "" {
		destPath = "Storage"
	}

	destPath = strings.TrimRight(destPath, "/")

	validatedPath, err := validatePathName(destPath, h.Config)
	if err != nil {
		http.Error(w, "Failed to validate pathname: "+err.Error(), http.StatusBadRequest)
		return
	}

	filename := r.FormValue("filename")
	if filename == "" {
		http.Error(w, "Filename not specified", http.StatusBadRequest)
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

	fullFilePath := filepath.Join(absolutePath, filename)

	dst, err := os.Create(fullFilePath)
	if err != nil {
		http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully as %s", fullFilePath)
}
