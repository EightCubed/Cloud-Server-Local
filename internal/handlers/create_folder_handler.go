package handlers

import "net/http"

func (h Handler) CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger
	logger.Info("Calling create folder handler")

}
