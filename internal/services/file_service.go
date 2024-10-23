package services

import (
	"io/fs"
	"os"
)

func ListFiles(path string) []fs.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		// logger.Fatal("Error!", zap.Any("err", err))
	}

	return entries
}
