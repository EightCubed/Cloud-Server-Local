package handlers

import (
	"fmt"
	"io/fs"
	"sort"

	"cloud-server/internal/models"
	"cloud-server/internal/services"

	"go.uber.org/zap"
)

// Main function to list files by depth and build the tree structure
func listFilesByDepthMain(fileName string, maxDepth int, logger *zap.Logger) *models.Node {
	if fileName == "" {
		fileName = "./uploads/"
	}
	fileTree := &models.Node{
		File: models.File{
			FileName: "Uploads",
			FileType: models.FileTypeFolder,
		},
		Adjacent: []*models.Node{},
	}

	logger.Info("Listing files main function")

	fileTree.Adjacent = listFilesByDepthRecursive(fileName, maxDepth, 0, logger)

	// fmt.Println("File Tree Structure:")
	// printFileTree(fileTree, 0)
	// fmt.Print("inner", fileTree)
	return fileTree
}

type ByType []fs.DirEntry

func (a ByType) Len() int      { return len(a) }
func (a ByType) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByType) Less(i, j int) bool {
	if a[i].IsDir() && !a[j].IsDir() {
		return true
	}
	if !a[i].IsDir() && a[j].IsDir() {
		return false
	}
	return a[i].Name() < a[j].Name()
}

// Recursive function to build the file tree with a depth limit
func listFilesByDepthRecursive(pathName string, maxDepth, currentDepth int, logger *zap.Logger) []*models.Node {
	logger.Info("Entered listing recursive function function")
	if currentDepth > maxDepth {
		return []*models.Node{}
	}

	entries := services.ListFiles(pathName, logger)
	sort.Sort(ByType(entries))
	fileTree := []*models.Node{}

	for _, entry := range entries {
		fileName := entry.Name()
		fileInfo, err := entry.Info()
		isDir := fileInfo.IsDir()
		if err != nil {
			logger.Error("Error getting file info", zap.Error(err))
			continue
		}

		fileType := models.FileTypeFile
		if isDir {
			fileType = models.FileTypeFolder
		}

		newNode := &models.Node{
			File: models.File{
				FileName: fileName,
				FileType: fileType,
			},
			Adjacent: []*models.Node{},
		}

		if isDir && currentDepth < maxDepth {
			newNode.Adjacent = listFilesByDepthRecursive(pathName+"/"+fileName, maxDepth, currentDepth+1, logger)
		}
		fileTree = append(fileTree, newNode)
	}

	return fileTree
}

func PrintSpaces(spaces int) {
	for i := 0; i < spaces; i++ {
		fmt.Print("\t")
	}
}

func printFileTree(node *models.Node, depth int) {
	PrintSpaces(depth)
	fmt.Println(node.File)

	// Print all adjacent nodes (children)
	for _, child := range node.Adjacent {
		printFileTree(child, depth+1)
	}
}
