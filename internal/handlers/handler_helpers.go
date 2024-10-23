package handlers

import (
	"fmt"

	"cloud-server/internal/models"
	"cloud-server/internal/services"
)

// Main function to list files by depth and build the tree structure
func listFilesByDepthMain(fileName string, maxDepth int) *models.Node {
	if fileName == "" {
		fileName = "./uploads/"
	}
	fileTree := &models.Node{
		Data:     "Uploads",
		Adjacent: []*models.Node{},
	}

	fileTree.Adjacent = listFilesByDepthRecursive(fileName, maxDepth, 0)

	// fmt.Println("File Tree Structure:")
	// printFileTree(fileTree, 0)
	// fmt.Print("inner", fileTree)
	return fileTree
}

// Recursive function to build the file tree with a depth limit
func listFilesByDepthRecursive(pathName string, maxDepth, currentDepth int) []*models.Node {
	if currentDepth > maxDepth {
		return nil
	}

	entries := services.ListFiles(pathName)
	var fileTree []*models.Node

	for _, entry := range entries {
		fileName := entry.Name()
		fileInfo, err := entry.Info()
		isDir := fileInfo.IsDir()
		if err != nil {
			fmt.Printf("Error getting file info: %v\n", err)
			continue
		}

		newNode := &models.Node{
			Data:     fileName,
			Adjacent: []*models.Node{},
		}
		fileTree = append(fileTree, newNode)

		if isDir && currentDepth < maxDepth {
			newNode.Adjacent = listFilesByDepthRecursive(pathName+"/"+fileName, maxDepth, currentDepth+1)
		}
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
	fmt.Println(node.Data)

	// Print all adjacent nodes (children)
	for _, child := range node.Adjacent {
		printFileTree(child, depth+1)
	}
}
