package main

import (
	"fmt"
	"io"
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data to retrieve the file
	err := r.ParseMultipartForm(10 << 20) // Limit the size of the file to 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file on the server
	dst, err := os.Create("./uploads/uploaded_file")
	if err != nil {
		http.Error(w, "Unable to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the new file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Unable to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond to the client
	fmt.Fprintln(w, "File uploaded successfully!")
}

func main() {
	// Serve files at the "/download" endpoint

	os.MkdirAll("./uploads", os.ModePerm)

	// Serve the upload handler at the "/upload" endpoint
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/download", fileDownloadHandler)

	// Start the HTTP server on port 8080
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
