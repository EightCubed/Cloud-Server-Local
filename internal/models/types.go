package models

import (
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.SugaredLogger
	Config Config
}

type Config struct {
	PATH_TO_DIRECTORY string
	STORAGE_DIRECTORY string
}
