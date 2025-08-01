package models

type DirectoryRequestBody struct {
	Directory string `json:"directory"`
}

type BreadCrumbType struct {
	Title        string `json:"title"`
	RelativePath string `json:"relativePath"`
}

type ResponseDataFileDirectory struct {
	Data    *Node            `json:"data"`
	Path    []BreadCrumbType `json:"path"`
	Message string           `json:"message"`
}

type DeleteRequestBody struct {
	FilesToBeDeleted []Node `json:"filesToBeDeleted"`
}

type DeleteResponseBody struct {
	SuccessCount int    `json:"success_count"`
	FailureCount int    `json:"failure_count"`
	Message      string `json:"message"`
}
