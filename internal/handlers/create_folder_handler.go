package handlers

import (
	"cloud-server/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (h Handler) CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger
	logger.Info("Calling create folder handler")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var reqData models.DirectoryRequestBody
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	validatedPathName, err := validatePathName(reqData.Directory, h.Config)
	if err != nil {
		http.Error(w, "Failed to validate pathname: "+err.Error(), http.StatusBadRequest)
		return
	}

	absolutePath := returnAbsolutePath(validatedPathName, h.Config)

	err = os.Mkdir(absolutePath, 0750)
	if err != nil {
		http.Error(w, "Failed to create folder: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Created folder successfully")
}
