package handlers

import "net/http"

func (h Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger
	logger.Info("Calling upload handler")

}
