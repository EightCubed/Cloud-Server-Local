package services

import (
	"io/fs"
	"os"

	"go.uber.org/zap"
)

func ListFiles(path string, logger *zap.Logger) []fs.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		logger.Error("Error!", zap.Any("err", err))
	}

	return entries
}
