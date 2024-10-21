package main

import (
	"fmt"
	"net/http"
	"os"
)

func fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Define the file path
	filePath := "./test/test.png"

	// Check if the file exists
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set headers to prompt the download
	w.Header().Set("Content-Disposition", "attachment; filename=test.png")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", getFileSize(filePath)))

	// Serve the file to the client
	http.ServeFile(w, r, filePath)
}

// Function to get the file size
func getFileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

func main() {
	// Serve files at the "/download" endpoint
	http.HandleFunc("/download", fileDownloadHandler)

	// Start the HTTP server on port 8080
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
