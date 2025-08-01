package handlers

import (
	"cloud-server/internal/models"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.SugaredLogger
	Config *models.Config
}

func ReturnHandler(logger *zap.Logger, cfg *models.Config) *Handler {
	return &Handler{
		Logger: logger.Sugar(),
		Config: cfg,
	}
}

func validatePathName(pathName string, cfg *models.Config) (string, error) {
	if !strings.HasPrefix(pathName, cfg.STORAGE_DIRECTORY) {
		return "", fmt.Errorf("invalid pathname passed")
	}
	return pathName, nil
}

func returnAbsolutePath(pathName string, cfg *models.Config) string {
	return filepath.Join(cfg.PATH_TO_DIRECTORY, pathName)
}

func returnFormattedRelativePath(pathName string, cfg *models.Config) string {
	relativePath := pathName
	if strings.HasPrefix(relativePath, cfg.PATH_TO_DIRECTORY) {
		relativePath = strings.Replace(relativePath, cfg.PATH_TO_DIRECTORY, "", 1)
	}

	if !strings.HasSuffix(relativePath, "/") {
		relativePath += "/"
	}
	return relativePath
}

func buildBreadCrumbs(pathName string) (breadCrumbs []models.BreadCrumbType) {
	splitPathNames := strings.Split(pathName, "/")
	var relativePathName string
	for _, pathName := range splitPathNames {
		relativePathName = filepath.Join(relativePathName, pathName)
		item := &models.BreadCrumbType{
			Title:        pathName,
			RelativePath: relativePathName,
		}
		breadCrumbs = append(breadCrumbs, *item)
	}
	return
}

func returnParentDirectory(pathName string, cfg *models.Config) string {
	parentDir, _ := filepath.Split(pathName)
	return parentDir
}

func listDirectoryRecursive(cfg *models.Config, relativefilepath string, currentDepth, maxDepth int) (*models.Node, error) {
	if currentDepth > maxDepth {
		return nil, nil
	}

	formattedRelativePath := returnFormattedRelativePath(relativefilepath, cfg)
	absolutefilepath := returnAbsolutePath(relativefilepath, cfg)

	dirEntry, err := os.ReadDir(absolutefilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	rootNode := &models.Node{
		File: models.File{
			FileName:         filepath.Base(relativefilepath),
			FileType:         models.FileTypeFolder,
			AbsoluteFilePath: formattedRelativePath,
		},
		FilePath:        formattedRelativePath,
		ParentDirectory: returnParentDirectory(formattedRelativePath, cfg),
		Children:        []*models.Node{},
	}

	for _, entry := range dirEntry {
		entryRelPath := filepath.Join(relativefilepath, entry.Name())
		entryFormattedPath := filepath.Join(formattedRelativePath, entry.Name())
		entryAbsPath := filepath.Join(formattedRelativePath, entry.Name())

		node := &models.Node{
			File: models.File{
				FileName:         entry.Name(),
				AbsoluteFilePath: entryAbsPath,
			},
			FilePath:        entryFormattedPath,
			ParentDirectory: returnParentDirectory(formattedRelativePath, cfg),
			Children:        nil,
		}

		if entry.IsDir() {
			node.File.FileType = models.FileTypeFolder

			children, err := listDirectoryRecursive(cfg, entryRelPath, currentDepth+1, maxDepth)
			if err != nil {
				return nil, fmt.Errorf("failed to read directory recursively: %w", err)
			}
			if children != nil {
				node.Children = children.Children
			}
		} else {
			node.File.FileType = models.FileTypeFile
		}

		rootNode.Children = append(rootNode.Children, node)
	}

	return rootNode, nil
}
