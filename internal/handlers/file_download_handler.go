package handlers

import (
	"cloud-server/pkg/utils"
	"fmt"
	"net/http"
	"os"
)

func (h Handler) FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger
	logger.Info("Calling create folder handler")

	fileName := r.URL.Query().Get("fileName")
	validatedPathName, err := validatePathName(fileName, h.Config)
	if err != nil {
		http.Error(w, "Failed to validate pathname: "+err.Error(), http.StatusBadRequest)
		return
	}

	absolutePath := returnAbsolutePath(validatedPathName, h.Config)

	file, err := os.Open(absolutePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", utils.GetFileSize(absolutePath)))

	http.ServeFile(w, r, absolutePath)
}
