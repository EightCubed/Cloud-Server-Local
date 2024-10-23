package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
}

type DirectoryRequestBody struct {
	Directory string
}

type Node struct {
	Data     string
	Adjacent []*Node
}

func (s *Server) fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Incoming Download request")

	fileName := r.URL.Query().Get("fileName")
	filePath := "./uploads/" + fileName

	s.logger.Info("File name to download", zap.String("File:", fileName))
	// Check if the file exists
	file, err := os.Open(filePath)
	if err != nil {
		s.logger.Error("File not found", zap.String("File:", fileName))
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set headers to prompt the download
	w.Header().Set("Content-Disposition", "attachment; filename=sample.txt")
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

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Incoming Upload")
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file: "+err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Info("File details: ", zap.Any("header", header))
	defer file.Close()

	fileName := header.Filename

	// Create a new file on the server
	dst, err := os.Create("./uploads/" + fileName)
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

func (s *Server) createFolder(w http.ResponseWriter, r *http.Request) {
	var request DirectoryRequestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Unable to decode body data: "+err.Error(), http.StatusBadRequest)
	}
	err = os.Mkdir("./uploads/"+request.Directory, 0700)
	if err != nil {
		http.Error(w, "Unable to create directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) listFiles(path string) []fs.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		s.logger.Fatal("Error!", zap.Any("err", err))
	}

	return entries
}

func PrintSpaces(spaces int) {
	for i := 0; i < spaces; i++ {
		fmt.Print("\t")
	}
}

func (s *Server) listFileDirectory(w http.ResponseWriter, r *http.Request) {
	treeStructure := s.listFilesByDepthMain("./uploads/", 10)
	fmt.Print("main", treeStructure)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(treeStructure) // Encode the file tree into JSON and send it
	if err != nil {
		http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Main function to list files by depth and build the tree structure
func (s *Server) listFilesByDepthMain(fileName string, maxDepth int) *Node {
	if fileName == "" {
		fileName = "./uploads/"
	}
	fileTree := &Node{
		Data:     "Uploads",
		Adjacent: []*Node{},
	}

	// Build the file tree recursively
	fileTree.Adjacent = s.listFilesByDepthRecursive(fileName, maxDepth, 0)

	// Print the entire file tree after it's built
	fmt.Println("File Tree Structure:")
	s.printFileTree(fileTree, 0)
	fmt.Print("inner", fileTree)
	return fileTree
}

// Recursive function to build the file tree with a depth limit
func (s *Server) listFilesByDepthRecursive(pathName string, maxDepth, currentDepth int) []*Node {
	if currentDepth > maxDepth {
		return nil
	}

	entries := s.listFiles(pathName)
	var fileTree []*Node

	for _, entry := range entries {
		fileName := entry.Name()
		fileInfo, err := entry.Info()
		isDir := fileInfo.IsDir()
		if err != nil {
			fmt.Printf("Error getting file info: %v\n", err)
			continue
		}

		// fileType := "File"
		// if isDir {
		// 	fileType = "Folder"
		// }

		// PrintSpaces(currentDepth)
		// fmt.Printf("fileName: %v (%v)\n", fileName, fileType)

		newNode := &Node{
			Data:     fileName,
			Adjacent: []*Node{},
		}
		fileTree = append(fileTree, newNode)

		if isDir && currentDepth < maxDepth {
			newNode.Adjacent = s.listFilesByDepthRecursive(pathName+"/"+fileName, maxDepth, currentDepth+1)
		}
	}

	return fileTree
}

func (s *Server) printFileTree(node *Node, depth int) {
	PrintSpaces(depth)
	fmt.Println(node.Data)

	// Print all adjacent nodes (children)
	for _, child := range node.Adjacent {
		s.printFileTree(child, depth+1)
	}
}

func main() {
	r := mux.NewRouter()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	server := &Server{logger: logger}

	// Serve files at the "/download" endpoint
	os.MkdirAll("./uploads", os.ModePerm)

	// Serve the upload handler at the "/upload" endpoint
	r.HandleFunc("/upload", server.uploadHandler).Methods("GET")

	// r.HandleFunc("/upload/{id:[A-Za-z0-9-_]+}", server.uploadHandler).Methods("GET")

	r.HandleFunc("/download", server.fileDownloadHandler).Methods("GET")

	r.HandleFunc("/createDirectory", server.createFolder).Methods("GET")

	r.HandleFunc("/showTreeDirectory", server.listFileDirectory).Methods("GET")

	// server.listFiles("Folder 1")
	server.listFilesByDepthMain("./uploads/", 10)

	// Start the HTTP server on port 8080
	// fmt.Println("Server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
