package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

		pathName := r.FormValue("path")

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to retrieve file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileName := header.Filename
		if fileName == ".DS_Store" || strings.HasPrefix(fileName, "._") {
			http.Error(w, "System files like .DS_Store are not allowed", http.StatusBadRequest)
			return
		}

		fullFilePath := filepath.Join(pathName, fileName)

		if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
			http.Error(w, "Failed to create directory: "+err.Error(), http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(fullFilePath)
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

// UploadHandler handles file uploads
func DeleteHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Incoming Delete request")
		var request models.DeleteRequestBody
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Unable to decode body data: "+err.Error(), http.StatusBadRequest)
			return
		}

		var resp models.DeleteResponseBody
		resp.Message = "Ok! I deleted"

		for _, elementToBeDeleted := range request.FilesToBeDeleted {
			logger.Info("Deleting file", zap.String("absolutefilePath", elementToBeDeleted.File.AbsoluteFilePath))
			successCount, failureCount := deleteByDepthSearch(elementToBeDeleted)
			resp.SuccessCount += successCount
			resp.FailureCount += failureCount
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
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
		err = os.Mkdir(request.Directory, 0700)
		if err != nil {
			http.Error(w, "Unable to create directory: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ListFileDirectoryHandler lists files in the parent directory
func ListFileDirectoryHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Listing directory")
		treeStructure := listFilesByDepthMain("./uploads", 10, logger)

		var resp models.ResponseDataFileDirectory
		resp.Data = *treeStructure
		resp.Path = []models.BreadCrumbType{{Title: "uploads", AbsolutePath: "uploads"}}
		resp.Message = "Successful list"

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ListFilesByPath lists files in a particular directory
func ListFilesByPath(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Fix: Input sanitization before querying
		pathName := r.URL.Query().Get("fileName")
		if pathName == "" {
			pathName = "./uploads"
		}
		logger.Info("Listing directory")
		treeStructure := listFilesByDepthMain(pathName, 10, logger)

		logger.Info("Query params", zap.String("pathName", pathName))

		var resp models.ResponseDataFileDirectory
		resp.Data = *treeStructure

		breadCrumbList := delete_empty(strings.Split(strings.Replace(pathName, "./", "", -1), "/"))

		var breadCrumb []models.BreadCrumbType
		var absolutePath string
		for _, elem := range breadCrumbList {
			absolutePath += elem + "/"
			crumb := &models.BreadCrumbType{
				Title:        elem,
				AbsolutePath: absolutePath,
			}
			breadCrumb = append(breadCrumb, *crumb)
		}

		resp.Path = breadCrumb
		resp.Message = "Successful list"

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
