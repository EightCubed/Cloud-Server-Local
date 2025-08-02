package handlers

import (
	"cloud-server/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (h Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger
	logger.Info("Calling delete handler")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var reqData models.DeleteRequestBody
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	validatedPathName, err := validatePathName(reqData.File.FilePath, h.Config)
	if err != nil {
		http.Error(w, "Failed to validate pathname: "+err.Error(), http.StatusBadRequest)
		return
	}

	absolutePath := returnAbsolutePath(validatedPathName, h.Config)

	if reqData.File.File.FileType == models.FileTypeFile {
		err := deleteFile(absolutePath)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "File deleted successfully ")
	} else {
		err := deleteFolder(absolutePath)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "File deleted successfully ")
	}
}
