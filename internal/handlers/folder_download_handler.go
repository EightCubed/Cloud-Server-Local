package handlers

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func (h Handler) DownloadFolderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderID := vars["id"]

	if folderID == "" {
		http.Error(w, "Folder path parameter missing", http.StatusBadRequest)
		return
	}

	validatedPathName, err := validatePathName(folderID, h.Config)
	if err != nil {
		http.Error(w, "Invalid folder path: "+err.Error(), http.StatusBadRequest)
		return
	}

	absFolderPath := returnAbsolutePath(validatedPathName, h.Config)

	info, err := os.Stat(absFolderPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Folder not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error accessing folder: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if !info.IsDir() {
		http.Error(w, "Requested path is not a directory", http.StatusBadRequest)
		return
	}

	folderName := filepath.Base(absFolderPath)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+folderName+`.zip"`)

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	err = filepath.Walk(absFolderPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(absFolderPath, path)
		if err != nil {
			return err
		}

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		return err
	})

	if err != nil {
		http.Error(w, "Error creating zip archive: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
