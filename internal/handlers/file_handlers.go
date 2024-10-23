package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"

	"cloud-server/internal/models"
	"cloud-server/pkg/utils"
)

// UploadHandler handles file uploads
func UploadHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Incoming Upload")
		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to retrieve file: "+err.Error(), http.StatusBadRequest)
			return
		}
		logger.Info("File details: ", zap.Any("header", header))
		defer file.Close()

		fileName := header.Filename

		// Create a new file on the server
		dst, err := os.Create("./uploads/" + fileName)
		if err != nil {
			http.Error(w, "Unable to create file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Unable to save file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "File uploaded successfully!")
	}
}

// FileDownloadHandler handles file downloads
func FileDownloadHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Incoming Download request")
		fileName := r.URL.Query().Get("fileName")
		filePath := "./uploads/" + fileName

		logger.Info("File name to download", zap.String("File:", fileName))
		file, err := os.Open(filePath)
		if err != nil {
			logger.Error("File not found", zap.String("File:", fileName))
			http.Error(w, "File not found.", http.StatusNotFound)
			return
		}
		defer file.Close()

		// Set headers to prompt the download
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", utils.GetFileSize(filePath)))

		http.ServeFile(w, r, filePath)
	}
}

// CreateFolderHandler handles directory creation
func CreateFolderHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Incoming Create Directory request")
		var request models.DirectoryRequestBody
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Unable to decode body data: "+err.Error(), http.StatusBadRequest)
			return
		}
		err = os.Mkdir("./uploads/"+request.Directory, 0700)
		if err != nil {
			http.Error(w, "Unable to create directory: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ListFileDirectoryHandler lists files in the directory
func ListFileDirectoryHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Listing directory")
		treeStructure := listFilesByDepthMain("./uploads/", 10)

		var resp models.ResponseDataFileDirectory
		resp.Data = *treeStructure
		resp.Message = "Successful list"

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
