package handlers

import (
	"cloud-server/internal/models"
	"encoding/json"
	"net/http"
)

const MAX_DEPTH = 3

func (h Handler) ListDirectoryTreeHandler(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger
	logger.Info("Calling list files handler")

	pathName := r.URL.Query().Get("path")

	pathName, err := validatePathName(pathName, h.Config)
	if err != nil {
		http.Error(w, "Failed to validate pathname: "+err.Error(), http.StatusBadRequest)
		return
	}

	fileNodeList, err := listDirectoryRecursive(h.Config, pathName, 0, MAX_DEPTH)
	if err != nil {
		http.Error(w, "Failed to read from directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	breadCrumbs := buildBreadCrumbs(pathName)

	resp := &models.ResponseDataFileDirectory{
		Data:    fileNodeList,
		Path:    breadCrumbs,
		Message: "Successfully listed directory",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
