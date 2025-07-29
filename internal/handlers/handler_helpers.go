package handlers

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"cloud-server/internal/models"
	"cloud-server/internal/services"

	"go.uber.org/zap"
)

// Main function to list files by depth and build the tree structure
func listFilesByDepthMain(fileName string, maxDepth int, logger *zap.Logger) *models.Node {
	if fileName == "" {
		fileName = "./uploads"
	}
	fileTree := &models.Node{
		File: models.File{
			FileName:         fileName,
			FileType:         models.FileTypeFolder,
			AbsoluteFilePath: fileName,
		},
		Children:        []*models.Node{},
		FilePath:        fileName,
		ParentDirectory: "",
	}

	if fileName != "" {
		splitStr := strings.Split(strings.Replace(fileName, "./", "", -1), "/")
		fileTree.ParentDirectory = strings.Join(splitStr[:len(splitStr)-1], "/")
	}

	fileTree.Children = listFilesByDepthRecursive(fileName, maxDepth, 0, logger)

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
				FileName:         fileName,
				FileType:         fileType,
				AbsoluteFilePath: pathName + "/" + fileName,
			},
			Children: []*models.Node{},
			FilePath: pathName,
		}

		if fileName != "" {
			splitStr := strings.Split(strings.Replace(fileName, "./", "", -1), "/")
			newNode.ParentDirectory = strings.Join(splitStr[:len(splitStr)-1], "/")
		}

		if isDir && currentDepth < maxDepth {
			newNode.Children = listFilesByDepthRecursive(pathName+"/"+fileName, maxDepth, currentDepth+1, logger)
		}
		fileTree = append(fileTree, newNode)
	}

	return fileTree
}

func deleteByDepthSearch(toDelete models.Node) (successCount int, failureCount int) {
	if toDelete.File.FileType == models.FileTypeFile {
		if err := deleteFile(toDelete.File.AbsoluteFilePath); err != nil {
			failureCount++
		} else {
			successCount++
		}
		return
	}

	for _, child := range toDelete.Children {
		sc, fc := deleteByDepthSearch(*child)
		successCount += sc
		failureCount += fc

		if child.File.FileType == models.FileTypeFolder {
			if err := deleteFile(child.File.AbsoluteFilePath); err != nil {
				failureCount++
			} else {
				successCount++
			}
		}
	}

	if err := deleteFile(toDelete.File.AbsoluteFilePath); err != nil {
		failureCount++
	} else {
		successCount++
	}

	return
}

func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func PrintSpaces(spaces int) {
	for i := 0; i < spaces; i++ {
		fmt.Print("\t")
	}
}

func printFileTree(node *models.Node, depth int) {
	PrintSpaces(depth)
	fmt.Println(node.File)

	for _, child := range node.Children {
		printFileTree(child, depth+1)
	}
}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
